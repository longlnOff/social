package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/longlnOff/social/cmd/configuration"
	"github.com/longlnOff/social/internal/store"
)

type application struct {
	configuration configuration.Configuration
	store store.Storage
}





func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(middleware.Timeout(time.Second * 60))

	r.Route("/v1", func (r chi.Router) {
		r.Get("/health", app.healthcheckHandler)
	
		r.Route("/posts", func (r chi.Router) {
			r.Post("/", app.createPostHandler)
		
			r.Route("/{postID}", func (r chi.Router) {
				r.Use(app.postContextMiddleware)
			
				r.Get("/", app.getPostsHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)
				
				// Create comment for post
				r.Post("/comments", app.createCommentHandler)
			})
		})
	
	})
	
	return r
}

func (app *application) run(mux http.Handler) error {
	server := http.Server{
		Addr:    app.configuration.Server.SERVER_ADDRESS + ":" + app.configuration.Server.SERVER_PORT,
		Handler: mux,
	}
	log.Printf("Starting server on %s port", app.configuration.Server.SERVER_PORT)
	return server.ListenAndServe()
}
