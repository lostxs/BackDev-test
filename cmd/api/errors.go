package main

import (
	"log"
	"net/http"
)

func (a *app) badRequestException(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s %s: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (a *app) unauthorizedException(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s %s: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (a *app) internalServerException(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("%s %s: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusInternalServerError, err.Error())
}
