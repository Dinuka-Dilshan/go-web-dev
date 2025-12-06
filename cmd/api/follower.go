package main

import "net/http"

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	type Body struct {
		UserID int `json:"user_id"`
	}

	var body Body

	if err := readJson(w, r, &body); err != nil {
		app.badRequestError(w, r, err)
		return
	}

	posts, err := app.store.Posts.GetUserFeed(r.Context(), body.UserID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	app.jsonResponse(w, http.StatusOK, posts)

}
