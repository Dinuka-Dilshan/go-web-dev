package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserId    int       `json:"user_id"`
	Tags      []string  `json:"tags"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Comments  []Comment `json:"comments"`
	Version   int       `json:"version"`
}

type PostWithMetaData struct {
	Post
	CommentCount int `json:"comments_count,omitempty"`
	User         struct {
		ID       int    `json:"id"`
		UserName string `json:"user_name"`
	} `json:"user"`
}

type PostStore struct {
	db *pgxpool.Pool
}

func (postStore *PostStore) GetUserFeed(
	ctx context.Context,
	userId int,
	pagination PaginatedQuery,
) ([]*PostWithMetaData, error) {

	query := `
			SELECT 
				p.id, 
				p.title, 
				p.user_id,
				p.content,
				p.created_at, 
				p.tags, 
				COALESCE(comment_counts.total, 0) AS comments_count,
				u.username
			FROM posts p
			LEFT JOIN (
				SELECT post_id, COUNT(*) AS total
				FROM comments
				GROUP BY post_id
			) comment_counts ON comment_counts.post_id = p.id
			JOIN followers f ON f.follower_id = p.user_id AND f.user_id = $1
			LEFT JOIN users u ON u.id = p.user_id
			WHERE 
				(p.title ILIKE '%'|| $4 || '%' OR p.content ILIKE '%'|| $4 || '%') AND
				(p.tags @> $5 OR p.tags @> '{}') AND
				(p.created_at >= $6 OR $6 IS NULL) AND
				(p.created_at < $7 OR $7 IS NULL)
			ORDER BY p.created_at ` + pagination.Sort + `
			LIMIT $2 OFFSET $3
			`

	rows, err := postStore.db.Query(
		ctx,
		query,
		userId,
		pagination.Limit,
		pagination.Offset,
		pagination.Search,
		pagination.Tags,
		pagination.Since,
		pagination.Until,
	)
	if err != nil {
		return nil, err
	}

	posts, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (*PostWithMetaData, error) {
		var post PostWithMetaData
		if err := row.Scan(
			&post.ID,
			&post.Title,
			&post.UserId,
			&post.Content,
			&post.CreatedAt,
			&post.Tags,
			&post.CommentCount,
			&post.User.UserName,
		); err != nil {
			return nil, err
		}
		post.User.ID = post.UserId
		return &post, nil
	})
	if err != nil {
		return nil, err
	}

	return posts, nil

}

func (postStore *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content,title,user_id,tags) 
	          VALUES ($1,$2,$3,$4) RETURNING id, created_at, updated_at`

	err := postStore.db.QueryRow(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserId,
		post.Tags,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (postStore *PostStore) GetPostById(ctx context.Context, id int) (*Post, error) {
	query := `SELECT id, title, content, user_id, tags, created_at , version FROM posts WHERE id=$1`

	var post Post
	err := postStore.db.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserId,
		&post.Tags,
		&post.CreatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return &post, nil

}

func (postStore *PostStore) Delete(ctx context.Context, postId int) error {
	query := `DELETE FROM posts WHERE id = $1`

	cmd, err := postStore.db.Exec(
		ctx,
		query,
		postId,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return ErrorNotFound
	}

	return nil
}

func (postStore *PostStore) Update(ctx context.Context, post *Post) error {
	query := `UPDATE posts 
			  SET title = $1, content = $2, version = version + 1
			  WHERE id = $3 AND version = $4
			  RETURNING version`

	err := postStore.db.QueryRow(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return ErrorNotFound
		default:
			return err
		}

	}

	return nil
}
