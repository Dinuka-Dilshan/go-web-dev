package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type UserKey string

const userCtx UserKey = "userKey"

// GetUserHandler godoc
//
//	@Summary		Get user details
//	@Description	Returns user details by the provided id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int	true	"User ID"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{object}	map[string]string
//	@Router			/users/{userId} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) userContextMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "userId")
		userId, err := strconv.Atoi(param)

		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		user, err := app.store.Users.GetUserById(r.Context(), userId)
		if err != nil {
			switch err {
			case pgx.ErrNoRows:
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, userCtx, user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

type FollowUser struct {
	UserID int `json:"user_id"`
}

// FollowUserHandler godoc
//
//	@Summary		Follow a user
//	@Description	Follow a user by the provided user id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int			true	"User ID"
//	@Param			payload	body		FollowUser	true	"Follow payload"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		409		{string}	string	"Conflict"
//	@Failure		500		{object}	map[string]string
//	@Router			/users/{userId}/follow [post]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload FollowUser
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	followUser := getUserFromCtx(r)

	if err := app.store.Followers.Follow(r.Context(), followUser.ID, payload.UserID); err != nil {
		switch err {
		case store.ErrorConflict:
			app.conflictError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// UnfollowUserHandler godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user by the provided user id
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userId	path		int			true	"User ID"
//	@Param			payload	body		FollowUser	true	"Unfollow payload"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{object}	map[string]string
//	@Router			/users/{userId}/unfollow [post]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload FollowUser
	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	unfollowUser := getUserFromCtx(r)

	if err := app.store.Followers.Unfollow(r.Context(), unfollowUser.ID, payload.UserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if err := app.store.Users.Activate(r.Context(), token); err != nil {
		switch err {
		case store.ErrorNotFound:
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
