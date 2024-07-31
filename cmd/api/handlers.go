package main

// import (
// 	"backend/internal/models"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"strconv"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/golang-jwt/jwt/v4"
// )

// func (app *application) Home(w http.ResponseWriter, _ *http.Request) {
// 	var payload = struct {
// 		Status  string `json:"status"`
// 		Message string `json:"message"`
// 		Version string `json:"version"`
// 	}{
// 		Status:  "active",
// 		Message: "Go Movies up and running",
// 		Version: "1.0.0",
// 	}

// 	_ = app.writeJSON(w, http.StatusOK, payload)
// }