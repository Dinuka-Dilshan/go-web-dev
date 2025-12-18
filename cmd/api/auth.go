package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
	"github.com/google/uuid"
)

type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

type UserWithToken struct {
	store.User
	Token string `json:"token"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := getValidator().Struct(&payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user := store.User{
		UserName: payload.Username,
		Email:    payload.Email,
	}

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	plainToken := uuid.New().String()
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	if err := app.store.Users.CreateAndInvite(r.Context(), &user, hashToken, app.config.mail.exp); err != nil {
		switch err {
		case store.ErrDuplicateEmail:
			app.badRequestError(w, r, err)
			return
		case store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}

	}

	app.jsonResponse(w, http.StatusOK, UserWithToken{
		User:  user,
		Token: plainToken,
	})
}
