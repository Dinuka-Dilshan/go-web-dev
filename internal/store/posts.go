package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Post struct {
	ID        int      `json:"id`
	Content   string   `json:"content`
	Title     string   `json:"title`
	UserId    int      `json:"user_id`
	Tags      []string `json:"tags`
	UpdatedAt string   `json:"updated_at`
	CreatedAt string   `json:"created_at`
}

type PostStore struct {
	db *pgxpool.Pool
}

func (postStore *PostStore) Create(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (content,title,user_id,tags) 
	          VALUES ($1,$2,$3,$4) RETURNING id, created_at, updated_at`

	_, err := postStore.db.Exec(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserId,
		post.Tags,
	)

	if err != nil {
		return err
	}

	return nil
}
