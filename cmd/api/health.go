package main

import "net/http"

func (app *application) healthCheckHandler(writer http.ResponseWriter, res *http.Request) {
	writer.Write([]byte("ok"))
}
