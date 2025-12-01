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
}

type PostStore struct {
	db *pgxpool.Pool
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
	query := `SELECT id, title, content, user_id, tags, created_at FROM posts WHERE id=$1`

	var post Post
	err := postStore.db.QueryRow(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserId,
		&post.Tags,
		&post.CreatedAt)

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
			  SET title = $1, content = $2
			  WHERE id = $3`

	_, err := postStore.db.Exec(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
