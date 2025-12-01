package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Comment struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	PostId    int       `json:"post_id"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UserName  string    `json:"user_name"`
}

type CommentStore struct {
	db *pgxpool.Pool
}

func (commentStore *CommentStore) GetByPostId(ctx context.Context, postId int) (*[]Comment, error) {
	query := `SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, users.username FROM comments c
			  JOIN users ON users.id = c.user_id
  			  WHERE c.post_id = $1`

	rows, err := commentStore.db.Query(ctx, query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []Comment{}

	for rows.Next() {
		var comment Comment
		err := rows.Scan(&comment.ID, &comment.PostId, &comment.UserId, &comment.Content, &comment.CreatedAt, &comment.UserName)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return &comments, nil
}
