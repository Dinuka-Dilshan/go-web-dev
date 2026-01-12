package main

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/mailer"
	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
	"github.com/golang-jwt/jwt/v5"
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
		RoleId:   1,
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

		case store.ErrDuplicateUsername:
			app.badRequestError(w, r, err)

		default:
			app.internalServerError(w, r, err)

		}
		return
	}

	data := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.UserName,
		ActivationURL: plainToken,
	}

	err := app.mailer.Send(mailer.UserWelcomeTemplate, user.UserName, user.Email, data, false)

	if err != nil {
		if err := app.store.Users.Delete(r.Context(), user.ID); err != nil {
			app.internalServerError(w, r, err)
			return
		}
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, UserWithToken{
		User:  user,
		Token: plainToken,
	})
}

type CreateTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255"`
	Password string `json:"password" validate:"required,min=3,max=72"`
}

func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateTokenPayload

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := getValidator().Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	user, err := app.store.Users.GetUserByEmail(r.Context(), payload.Email)

	if err != nil {
		switch err {
		case store.ErrorNotFound:
			app.unauthorizedError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := user.Password.Compare(payload.Password); err != nil {
		app.unauthorizedError(w, r, err)
		return
	}

	claims := jwt.MapClaims{
		"iss": app.config.auth.iss,
		"aud": app.config.auth.aud,
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
	}

	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	if err := app.jsonResponse(w, http.StatusCreated, response{Token: token}); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
