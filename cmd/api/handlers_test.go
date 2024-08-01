package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/jwtauth"
	"github.com/sirupsen/logrus"
)

// MockAuthenticator is a mock implementation of the Authenticator interface
type MockAuthenticator struct {
	ValidateCredentialsFunc func(user, pass string) bool
	GenerateTokenFunc       func(user string) (string, error)
	TokenAuthInstance       *jwtauth.JWTAuth
}

func (m *MockAuthenticator) ValidateUserCredentials(user, pass string) bool {
	return m.ValidateCredentialsFunc(user, pass)
}

func (m *MockAuthenticator) GenerateToken(user string) (string, error) {
	return m.GenerateTokenFunc(user)
}

func (m *MockAuthenticator) TokenAuth() *jwtauth.JWTAuth {
	return m.TokenAuthInstance
}

func TestAuthenticateValidCredentials(t *testing.T) {
	authenticator := &MockAuthenticator{
		ValidateCredentialsFunc: func(user, pass string) bool {
			return user == "admin" && pass == "password"
		},
		GenerateTokenFunc: func(user string) (string, error) {
			return "mockToken", nil
		},
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	app := &application{
		authenticator: authenticator,
		logger:        logger,
	}

	reqBody := bytes.NewBufferString(`{"user":"admin","pass":"password"}`)
	req, err := http.NewRequest(http.MethodPost, "/authenticate", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.authenticate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	cookie := rr.Result().Cookies()[0]
	if cookie.Name != "jwtToken" || cookie.Value != "mockToken" {
		t.Errorf("handler did not set the correct cookie: got %v want %v", cookie.Value, "mockToken")
	}

	expected := `{"token":"mockToken"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAuthenticateInvalidCredentials(t *testing.T) {
	authenticator := &MockAuthenticator{
		ValidateCredentialsFunc: func(user, pass string) bool {
			return user == "admin" && pass == "password"
		},
		GenerateTokenFunc: func(user string) (string, error) {
			return "mockToken", nil
		},
	}

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	app := &application{
		authenticator: authenticator,
		logger:        logger,
	}

	reqBody := bytes.NewBufferString(`{"user":"admin","pass":"wrongpassword"}`)
	req, err := http.NewRequest(http.MethodPost, "/authenticate", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.authenticate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnauthorized)
	}

	expected := `{"error":true,"message":"Unauthorized"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
