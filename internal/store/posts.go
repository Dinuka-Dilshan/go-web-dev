package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        int       `json:"id`
	Content   string    `json:"content`
	Title     string    `json:"title`
	UserId    int       `json:"user_id`
	Tags      []string  `json:"tags`
	UpdatedAt time.Time `json:"updated_at`
	CreatedAt time.Time `json:"created_at`
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
