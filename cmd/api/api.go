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
	Server ServerConfiguration
	Database DatabaseConfiguration
}

type ServerConfiguration struct {
	SERVER_ADDRESS string	`mapstructure:"SERVER_ADDRESS"`
	SERVER_PORT string		`mapstructure:"SERVER_PORT"`
	ENVIRONMENT string		`mapstructure:"ENVIRONMENT"`
	VERSION string			`mapstructure:"VERSION"`
}

type DatabaseConfiguration struct {
	ENGINE					string			`mapstructure:"DB_ENGINE"`
	HOST 					string			`mapstructure:"DB_HOST"`
	PORT 					string			`mapstructure:"DB_PORT"`
	USER 					string			`mapstructure:"DB_USER"`
	PASSWORD 				string			`mapstructure:"DB_PASSWORD"`
	DB_NAME 				string			`mapstructure:"DB_NAME"`
	DB_MAX_OPEN_CONNS 		int				`mapstructure:"DB_MAX_OPEN_CONNS"`
	DB_MAX_IDLE_CONNS 		int				`mapstructure:"DB_MAX_IDLE_CONNS"`
	DB_MAX_IDLE_TIME 		time.Duration	`mapstructure:"DB_MAX_IDLE_TIME"`
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
		Addr:    app.configuration.Server.SERVER_ADDRESS + ":" + app.configuration.Server.SERVER_PORT,
		Handler: mux,
	}
	log.Printf("Starting server on %s port", app.configuration.Server.SERVER_PORT)
	return server.ListenAndServe()
}
