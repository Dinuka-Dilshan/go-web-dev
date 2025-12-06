package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Follower struct {
	UserId     string    `json:"user_id"`
	FollowerId string    `json:"follower_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type FollowerStore struct {
	db *pgxpool.Pool
}

func (followerStore *FollowerStore) Follow(ctx context.Context, followerID int, userId int) error {
	query := `INSERT INTO followers (user_id,follower_id)
			VALUES ($1,$2)
	`
	_, err := followerStore.db.Exec(ctx, query, userId, followerID)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return ErrorConflict
		}
		return err
	}

	return err
}

func (followerStore *FollowerStore) Unfollow(ctx context.Context, followerID int, userId int) error {
	query := `DELETE FROM followers 
			  WHERE user_id = $1 AND follower_id = $2
	`
	_, err := followerStore.db.Exec(ctx, query, userId, followerID)

	return err
}
