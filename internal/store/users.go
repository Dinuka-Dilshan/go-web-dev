package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int       `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type UserStore struct {
	db *pgxpool.Pool
}

func (userStore *UserStore) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users (username,password,email) 
			  VALUES ($1,$2,$3)
			  RETURNING id, created_at`

	err := userStore.db.QueryRow(
		ctx,
		query,
		user.UserName,
		user.Password,
		user.Email,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}

func (usersStore *UserStore) GetUserById(ctx context.Context, userId int) (*User, error) {

	query := `SELECT id, email, username,  created_at
			  FROM users
			  WHERE id=$1
			`
	var user = &User{}

	err := usersStore.db.QueryRow(
		ctx,
		query,
		userId,
	).Scan(
		&user.ID,
		&user.Email,
		&user.UserName,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
