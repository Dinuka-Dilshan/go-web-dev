package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int    `json:"id`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at`
}

type UserStore struct {
	db *pgxpool.Pool
}

func (userStore *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username,password,email) VALUES ($1,$2,$3)`

	_, err := userStore.db.Exec(
		ctx,
		query,
		user.UserName,
		user.Password,
		user.Email,
	)

	if err != nil {
		return err
	}

	return nil
}
