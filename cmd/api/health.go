package main

import (
	"log"
	"net/http"
)

func (app *application) healthCheckHandler(writer http.ResponseWriter, res *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"version": version,
	}

	err := writeJson(writer, 200, data)
	if err != nil {
		log.Print(err)
	}
}
