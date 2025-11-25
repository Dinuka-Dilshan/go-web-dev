package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const version = "0.0.1"

type application struct {
	config config
	store  store.Storage
}

type config struct {
	address  string
	dbConfig dbConfig
}

type dbConfig struct {
	address            string
	maxOpenConnections int32
	maxIdleTime        time.Duration
}

func (app *application) mount() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/post", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postId}", func(r chi.Router) {
				r.Get("/", app.getPostHandler)
			})
		})
	})

	return router
}

func (app *application) run(mux *http.Handler) error {

	server := &http.Server{
		Addr:         app.config.address,
		Handler:      *mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server is listning on %v", app.config.address)

	return server.ListenAndServe()
}
