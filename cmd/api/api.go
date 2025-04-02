package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/longlnOff/social/cmd/configuration"
	"github.com/longlnOff/social/docs"
	"github.com/longlnOff/social/internal/store"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"go.uber.org/zap"
)

type application struct {
	configuration configuration.Configuration
	store         store.Storage
	logger        *zap.Logger
}

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(middleware.Timeout(time.Second * 60))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthcheckHandler)

		docsURL := fmt.Sprintf("%s/swagger/doc.json",
			app.configuration.Server.SERVER_PORT)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL((docsURL))))
		// Post API
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)

			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postContextMiddleware)

				r.Get("/", app.getPostsHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)

				// Create comment for post
				r.Post("/comments", app.createCommentHandler)
			})

			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		// User API
		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)

				// Follow & Unfollow
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
		})

	})

	return r
}

func (app *application) run(mux http.Handler) error {

	// Setting docs swagger
	docs.SwaggerInfo.Version = app.configuration.Server.VERSION
	docs.SwaggerInfo.Host = app.configuration.Server.EXTERNAL_ADDRESS + ":" + app.configuration.Server.EXTERNAL_PORT
	docs.SwaggerInfo.BasePath = "/v1"

	server := http.Server{
		Addr:    app.configuration.Server.SERVER_ADDRESS + ":" + app.configuration.Server.SERVER_PORT,
		Handler: mux,
	}
	app.logger.Info("Starting server on:", zap.String("port", app.configuration.Server.SERVER_PORT))
	return server.ListenAndServe()
}
