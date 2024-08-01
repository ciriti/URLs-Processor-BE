package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)
	mux.Get("/test", app.regularEndpoint)

	mux.Post("/authenticate", app.authenticate)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(jwtauth.Verifier(app.authenticator.TokenAuth()))
		mux.Use(jwtauth.Authenticator)
		mux.Get("/test-protected", app.adminEndpoint)
	})

	return (mux)
}
