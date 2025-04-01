package main

import (
	"net/http"
)

type HealthResponse struct {
	Status      string `json:"status"`
	Environment string `json:"environment"`
	Version     string `json:"version"`
}

// Healthcheck godoc

// @Summary		Healthcheck
// @Description	Healthcheck
// @Tags			healthcheck
// @Produce		json
// @Success		200	{object}	HealthResponse
// @Failure		500	{object}	error
// @Router			/health [get]
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := HealthResponse{
		Status:      "available",
		Environment: app.configuration.Server.ENVIRONMENT,
		Version:     app.configuration.Server.VERSION,
	}
	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
