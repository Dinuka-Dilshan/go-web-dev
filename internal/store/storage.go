package store

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrorNotFound = errors.New("resource not found")
	ErrorConflict = errors.New("conflict")
)

type Storage struct {
	Posts interface {
		Create(context.Context, *Post) error
		GetPostById(context.Context, int) (*Post, error)
		Delete(context.Context, int) error
		Update(context.Context, *Post) error
	}

	Users interface {
		Create(context.Context, *User) error
		GetUserById(context.Context, int) (*User, error)
	}

	Comments interface {
		GetByPostId(context.Context, int) (*[]Comment, error)
		Create(context.Context, *Comment) error
	}

	Followers interface {
		Follow(ctx context.Context, followerId int, userId int) error
		Unfollow(ctx context.Context, followerId int, userId int) error
	}
}

func NewStorage(db *pgxpool.Pool) *Storage {
	return &Storage{
		Posts:     &PostStore{db},
		Users:     &UserStore{db},
		Comments:  &CommentStore{db},
		Followers: &FollowerStore{db},
	}
}
