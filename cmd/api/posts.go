package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	err := readJson(w, r, &payload)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post := &store.Post{
		Content: payload.Content,
		Title:   payload.Title,
		Tags:    payload.Tags,
		UserId:  1,
	}

	err = app.store.Posts.Create(r.Context(), post)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	writeJson(w, http.StatusCreated, post)

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "postId")
	postId, err := strconv.Atoi(param)

	if err != nil{
		app.badRequestError(w,r,err)
		return
	}

	post, err := app.store.Posts.GetPostById(r.Context(), postId)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	writeJson(w, http.StatusOK, post)
}
