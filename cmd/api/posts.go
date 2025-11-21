package main

import (
	"fmt"
	"net/http"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
)

type CreatePostPayload struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:tags`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	err := readJson(w, r, &payload)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	post := &store.Post{
		Content: payload.Content,
		Title:   payload.Title,
		Tags:    payload.Tags,
		UserId:  1,
	}

	if err = app.store.Posts.Create(r.Context(), post); err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		fmt.Print(err)
		return
	}

	writeJson(w, 200, post)

}
