package main

import (
	"net/http"
)

// HealthCheckHandler godoc
//
//	@Summary		Gopher Health check endpoint
//	@Description	Returns the health status of the API
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/health [get]
func (app *application) healthCheckHandler(writer http.ResponseWriter, res *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"version": version,
	}

	err := writeJson(writer, 200, data)
	if err != nil {
		app.logger.Error(err)
	}
}
