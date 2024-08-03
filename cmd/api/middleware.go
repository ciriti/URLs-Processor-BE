package main

import (
	"net/http"
)

func (app *application) enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowedOrigin := getEnv("ALLOWED_ORIGIN", "http://allowed-origin.com")

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {

			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
			return
		}
		h.ServeHTTP(w, r)
	})
}
