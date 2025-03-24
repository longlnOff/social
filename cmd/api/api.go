package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/longlnOff/social/internal/store"
)

type application struct {
	configuration configuration
	store store.Storage
}


type configuration struct {
	SERVER_ADDRESS string	`mapstructure:"SERVER_ADDRESS"`
	SERVER_PORT string		`mapstructure:"SERVER_PORT"`
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
	})

	
	return r
}




func (app *application) run(mux http.Handler) error {
	server := http.Server{
		Addr:    app.configuration.SERVER_ADDRESS + ":" + app.configuration.SERVER_PORT,
		Handler: mux,
	}
	log.Printf("Starting server on %s port", app.configuration.SERVER_PORT)
	return server.ListenAndServe()
}
