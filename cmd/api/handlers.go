package main

import (
	// "encoding/json"
	// "errors"
	// "fmt"
	// "io"
	// "log"
	"errors"
	"net/http"
	// "net/url"
	// "strconv"
	// "time"
	// "github.com/go-chi/chi/v5"
	// "github.com/golang-jwt/jwt/v4"
)

func (app *application) Home(w http.ResponseWriter, _ *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "URLs Processor",
		Version: "1.0.0",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	user := r.FormValue("user")
	pass := r.FormValue("pass")

	if user == "" || pass == "" {
		err := app.errorJSON(w, errors.New("username and password are required"), http.StatusBadRequest)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
		return
	}

	if !app.authenticator.ValidateUserCredentials(user, pass) {
		err := app.errorJSON(w, errors.New(http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
		return
	}

	tokenString, err := app.authenticator.GenerateToken(user)
	if err != nil {
		app.logger.WithError(err).Error("error generating token")
		err := app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, map[string]string{"token": tokenString}); err != nil {
		app.logger.WithError(err).Error("error writing JSON response")
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
	}
}
