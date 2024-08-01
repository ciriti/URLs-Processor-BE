package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockAuthenticator is a mock implementation of the Authenticator interface
type MockAuthenticator struct {
	ValidateCredentialsFunc func(user, pass string) bool
	GenerateTokenFunc       func(user string) (string, error)
}

func (m *MockAuthenticator) ValidateUserCredentials(user, pass string) bool {
	return m.ValidateCredentialsFunc(user, pass)
}

func (m *MockAuthenticator) GenerateToken(user string) (string, error) {
	return m.GenerateTokenFunc(user)
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

	app := &application{
		authenticator: authenticator,
	}

	reqBody := bytes.NewBufferString("user=admin&pass=password")
	req, err := http.NewRequest(http.MethodPost, "/authenticate", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.authenticate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `{"token":"mockToken"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestAuthenticateMissingCredentials(t *testing.T) {
	authenticator := &MockAuthenticator{
		ValidateCredentialsFunc: func(user, pass string) bool {
			return user == "admin" && pass == "password"
		},
		GenerateTokenFunc: func(user string) (string, error) {
			return "mockToken", nil
		},
	}

	app := &application{
		authenticator: authenticator,
	}

	reqBody := bytes.NewBufferString("user=&pass=")
	req, err := http.NewRequest(http.MethodPost, "/authenticate", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.authenticate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}

	expected := `{"error":true,"message":"username and password are required"}`
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

	app := &application{
		authenticator: authenticator,
	}

	reqBody := bytes.NewBufferString("user=admin&pass=wrongpassword")
	req, err := http.NewRequest(http.MethodPost, "/authenticate", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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
