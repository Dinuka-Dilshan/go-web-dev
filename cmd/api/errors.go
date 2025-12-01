package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error %s path: %s error:%s", r.Method, r.URL.Path, err.Error())

	writeJson(w, http.StatusInternalServerError, map[string]string{
		"error": "server encountered a problem",
	})
}

func (app *application) notFoundError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error %s path: %s error:%s", r.Method, r.URL.Path, err.Error())

	writeJson(w, http.StatusNotFound, map[string]string{
		"error": "not found",
	})
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error %s path: %s error:%s", r.Method, r.URL.Path, err.Error())

	writeJson(w, http.StatusBadRequest, map[string]string{
		"error": err.Error(),
	})
}
