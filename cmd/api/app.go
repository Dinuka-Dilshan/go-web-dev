package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dinuka-Dilshan/go-web-dev/docs"
	"github.com/Dinuka-Dilshan/go-web-dev/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

const version = "0.0.2"

type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	address  string
	dbConfig dbConfig
	apiUrl   string
	mail     mailConfig
}

type mailConfig struct {
	exp time.Duration
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
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("%v/v1/swagger/doc.json", app.config.address)),
		))

		r.Route("/post", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postId}", func(r chi.Router) {
				r.Use(app.postMiddleware)

				r.Delete("/", app.deletePostHandler)
				r.Get("/", app.getPostHandler)
				r.Patch("/", app.updatePostHandler)
			})

		})
		r.Route("/users", func(r chi.Router) {
			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.userContextMiddleWare)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
		})
	})

	return router
}

func (app *application) run(mux *http.Handler) error {

	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiUrl
	docs.SwaggerInfo.BasePath = "/v1"

	server := &http.Server{
		Addr:         app.config.address,
		Handler:      *mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("server is listning on port", "addr", app.config.address)

	return server.ListenAndServe()
}
