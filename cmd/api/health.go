package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "available",
		"environment": app.configuration.Server.ENVIRONMENT,
		"version": app.configuration.Server.VERSION,
	}
	if err := writeJSON(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err) 
	}
}
