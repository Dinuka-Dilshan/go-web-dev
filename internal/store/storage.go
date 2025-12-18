package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
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
		GetUserFeed(context.Context, int, PaginatedQuery) ([]*PostWithMetaData, error)
	}

	Users interface {
		Create(ctx context.Context, txn pgx.Tx, user *User) error
		GetUserById(context.Context, int) (*User, error)
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		Activate(context.Context, string) error
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

func withTransaction(db *pgxpool.Pool, ctx context.Context, fn func(pgx.Tx) error) error {
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}
