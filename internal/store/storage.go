package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrorNotFound = errors.New("resource not found")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetPostById(context.Context, int) (*Post, error)
	}

	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{
		Posts: &PostStore{db},
		Users: &UsersStore{db},
	}
}
