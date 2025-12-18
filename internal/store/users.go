package store

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type User struct {
	ID        int       `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	IsActive  bool      `json:"is_active"`
}

type password struct {
	text *string
	hash []byte
}

func (password *password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	password.hash = hash
	password.text = &text

	return nil
}

type UserStore struct {
	db *pgxpool.Pool
}

func (userStore *UserStore) Activate(ctx context.Context, token string) error {

	return withTransaction(userStore.db, ctx, func(tx pgx.Tx) error {
		user, err := userStore.getUserFromInvitation(ctx, tx, token)

		if err != nil {
			return err
		}

		user.IsActive = true
		if err := userStore.update(ctx, tx, user); err != nil {
			return err
		}

		if err := userStore.deleteUserInvitation(ctx, tx, user.ID); err != nil {
			return err
		}

		return nil

	})

}

func (useStore *UserStore) getUserFromInvitation(
	context context.Context,
	transaction pgx.Tx,
	token string,
) (*User, error) {

	query := `
		SELECT u.id, u.username, u.email, u.is_active
		FROM users u
		INNER JOIN user_invitations ui
		ON u.id = ui.user_id
		WHERE ui.token = $1 AND ui.expiry > $2
	`
	hash := sha256.Sum256([]byte(token))
	hashToken := hex.EncodeToString(hash[:])

	var user User
	fmt.Print(user.ID, hashToken)

	if err := transaction.QueryRow(context, query, hashToken, time.Now()).Scan(
		&user.ID,
		&user.UserName,
		&user.Email,
		&user.IsActive,
	); err != nil {
		switch err {
		case pgx.ErrNoRows:
			return nil, ErrorNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func (userStore *UserStore) Create(ctx context.Context, txn pgx.Tx, user *User) error {
	query := `INSERT INTO users (username,password,email) 
			  VALUES ($1,$2,$3)
			  RETURNING id, created_at`

	err := txn.QueryRow(
		ctx,
		query,
		user.UserName,
		user.Password.hash,
		user.Email,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_email_key\""):
			return ErrDuplicateEmail
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint \"users_username_key\""):
			return ErrDuplicateUsername
		default:
			return err
		}
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

func (usersStore *UserStore) CreateAndInvite(
	ctx context.Context,
	user *User,
	toekn string,
	invitaionExp time.Duration) error {

	return withTransaction(usersStore.db, ctx, func(tx pgx.Tx) error {
		if err := usersStore.Create(ctx, tx, user); err != nil {
			return err
		}

		if err := usersStore.createUserInvitation(ctx, tx, toekn, invitaionExp, user.ID); err != nil {
			return err
		}

		return nil
	})
}

func (userStore *UserStore) createUserInvitation(
	ctx context.Context,
	tx pgx.Tx,
	toekn string,
	invitaionExp time.Duration,
	userID int) error {
	query := `INSERT INTO user_invitations (token, user_id, expiry) VALUES ($1,$2,$3)`

	if _, err := tx.Exec(ctx, query, toekn, userID, time.Now().Add(invitaionExp)); err != nil {
		return err
	}

	return nil
}

func (userStore *UserStore) update(context context.Context, transaction pgx.Tx, user *User) error {
	query := `UPDATE users 
			SET username = $1, email = $2, is_active = $3
			WHERE users.id = $4
	`

	if _, err := transaction.Exec(
		context,
		query,
		user.UserName,
		user.Email,
		user.IsActive,
		user.ID,
	); err != nil {
		return err
	}

	return nil
}

func (userStore *UserStore) deleteUserInvitation(
	context context.Context,
	transaction pgx.Tx,
	userId int,
) error {
	query := `DELETE FROM user_invitations WHERE user_id =$1`

	if _, err := transaction.Exec(context, query, userId); err != nil {
		return err
	}

	return nil
}
