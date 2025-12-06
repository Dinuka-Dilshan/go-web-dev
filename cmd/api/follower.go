package main

import (
	"net/http"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	var paginatedQuery = store.PaginatedQuery{
		Limit:  10,
		Offset: 0,
		Sort:   "DESC",
	}

	if err := paginatedQuery.Parse(r); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	if err := getValidator().Struct(paginatedQuery); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	type Body struct {
		UserID int `json:"user_id"`
	}

	var body Body

	if err := readJson(w, r, &body); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	posts, err := app.store.Posts.GetUserFeed(
		r.Context(),
		body.UserID,
		paginatedQuery,
	)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, posts)

}
