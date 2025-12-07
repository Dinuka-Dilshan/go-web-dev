package main

import (
	"net/http"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
)

// GetUserFeedHandler godoc
//
//	@Summary		Get user feed
//	@Description	Returns paginated feed of posts for a user
//	@Tags			posts
//	@Accept			json
//	@Produce		json
//	@Param			limit	query		int		false	"Limit"						default(10)
//	@Param			offset	query		int		false	"Offset"					default(0)
//	@Param			sort	query		string	false	"Sort order (ASC or DESC)"	default(DESC)
//	@Param			body	body		object	true	"User ID"
//	@Success		200		{array}		map[string]interface{}
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{object}	map[string]string
//	@Router			/post/feed [post]
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
