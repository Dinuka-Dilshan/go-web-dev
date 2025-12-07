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

// CreatePostHandler godoc
//
//	@Summary		Create a new post
//	@Description	Creates a new post with the provided title, content, and tags
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePostPayload	true	"Post payload"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{object}	map[string]string
//	@Router			/post [post]
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

	app.jsonResponse(w, http.StatusCreated, post)

}

// GetPostHandler godoc
//
//	@Summary		Get post details
//	@Description	Returns post details by the provided id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postId	path		int	true	"Post ID"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{object}	map[string]string
//	@Router			/post/{postId} [get]
func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostId(r.Context(), post.ID)

	if err != nil {
		log.Printf("internal server error %s path: %s error:%s", r.Method, r.URL.Path, err.Error())
	} else {
		post.Comments = *comments
	}

	app.jsonResponse(w, http.StatusOK, post)
}

// DeletePostHandler godoc
//
//	@Summary		Delete post
//	@Description	Delete post details by the provided id
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			postId	path		int	true	"Post ID"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{object}	map[string]string
//	@Router			/post/{postId} [delete]
func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	err := app.store.Posts.Delete(r.Context(), post.ID)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrorNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	app.jsonResponse(w, http.StatusOK, nil)
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {

	var payload struct {
		Title   *string `json:"title" validate:"omitempty,max=100"`
		Content *string `json:"content" validate:"omitempty,max=1000"`
	}

	// UpdatePostHandler godoc
	//
	//	@Summary		Update post
	//	@Description	Update post details by the provided id
	//	@Tags			posts
	//	@Accept			json
	//	@Produce		json
	//	@Param			postId	path		int		true	"Post ID"
	//	@Param			payload	body		object	true	"Update post payload"
	//	@Success		201		{object}	map[string]interface{}
	//	@Failure		400		{string}	string	"Bad request"
	//	@Failure		500		{object}	map[string]string
	//	@Router			/post/{postId} [patch]

	if err := readJson(w, r, &payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := getValidator().Struct(payload); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	post := getPostFromCtx(r)

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

	app.jsonResponse(w, http.StatusCreated, nil)

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

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postCtx).(*store.Post)
	return post
}
