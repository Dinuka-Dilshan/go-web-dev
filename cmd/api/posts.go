package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
	"github.com/go-chi/chi/v5"
)

type postKey string

const postCtx postKey = "postKey"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
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

	if err = getValidator().Struct(payload); err != nil {
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
	post := getPostFromCtx(r.Context())

	comments, err := app.store.Comments.GetByPostId(r.Context(), post.ID)

	if err != nil {
		log.Printf("internal server error %s path: %s error:%s", r.Method, r.URL.Path, err.Error())
	} else {
		post.Comments = *comments
	}

	writeJson(w, http.StatusOK, post)
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "postId")
	postId, err := strconv.Atoi(param)

	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

	err = app.store.Posts.Delete(r.Context(), postId)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	writeJson(w, http.StatusOK, nil)
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {

	var payload struct {
		Title   *string `json:"title" validate:"omitempty,max=100"`
		Content *string `json:"content" validate:"omitempty,max=1000"`
	}

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := getValidator().Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post := getPostFromCtx(r.Context())

	if payload.Content != nil {
		post.Content = *payload.Content
	}
	if payload.Title != nil {
		post.Title = *payload.Title
	}

	err := app.store.Posts.Update(r.Context(), post)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	writeJson(w, http.StatusCreated, nil)

}

func (app *application) postMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		param := chi.URLParam(r, "postId")
		postId, err := strconv.Atoi(param)

		if err != nil {
			app.badRequestError(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetPostById(ctx, postId)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrorNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postCtx, post)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

func getPostFromCtx(ctx context.Context) *store.Post {
	post, _ := ctx.Value(postCtx).(*store.Post)
	return post
}
