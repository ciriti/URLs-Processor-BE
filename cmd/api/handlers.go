package main

import (
	// "encoding/json"
	// "errors"
	// "fmt"
	// "io"
	// "log"
	"errors"
	"log"
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
            log.Printf("error writing JSON response: %v", err)
        }
        return
    }

	if !app.authenticator.ValidateUserCredentials(user, pass) {
		err := app.errorJSON(w, errors.New(http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
		if err != nil {
			log.Printf("error writing JSON response: %v", err)
		}
		return
	}

	tokenString, err := app.authenticator.GenerateToken(user)
	if err != nil {
		log.Printf("error generating token: %v", err)
		err := app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			log.Printf("error writing JSON response: %v", err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, map[string]string{"token": tokenString}); err != nil {
		log.Printf("error writing JSON response: %v", err)
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			log.Printf("error writing JSON response: %v", err)
		}
	}
}
