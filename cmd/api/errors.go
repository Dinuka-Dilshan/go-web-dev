package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJson(w, http.StatusInternalServerError, map[string]string{
		"error": "server encountered a problem",
	})
}

func (app *application) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJson(w, http.StatusConflict, map[string]string{
		"error": "resource already exsists",
	})
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJson(w, http.StatusNotFound, map[string]string{
		"error": "not found",
	})
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("bad request error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeJson(w, http.StatusBadRequest, map[string]string{
		"error": err.Error(),
	})
}
