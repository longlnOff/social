package main

import (
	"go.uber.org/zap"
	"net/http"
)

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("Conflict:", zap.String("error", err.Error()), zap.String("path", r.URL.Path), zap.String("method", r.Method))
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("Internal server error:", zap.String("error", err.Error()), zap.String("path", r.URL.Path), zap.String("method", r.Method))
	writeJSONError(w, http.StatusInternalServerError, "Internal server error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("Bad request:", zap.String("error", err.Error()), zap.String("path", r.URL.Path), zap.String("method", r.Method))
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("Resource not found:", zap.String("error", err.Error()), zap.String("path", r.URL.Path), zap.String("method", r.Method))
	writeJSONError(w, http.StatusNotFound, "Resource not found")
}
