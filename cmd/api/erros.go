package main

import (
	"log"
	"net/http"
)

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Conflict: %v path: %v error: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: %v path: %v error: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusInternalServerError, "Internal server error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request: %v path: %v error: %v", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request: %v path: %v error: %v", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusNotFound, "Resource not found")
}
