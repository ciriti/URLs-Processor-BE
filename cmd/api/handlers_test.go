package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/jwtauth"
	"github.com/sirupsen/logrus"
)

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

type MockTaskQueue struct {
	AddTaskFunc  func(urlInfo *URLInfo) (*Task, error)
	StopTaskFunc func(id int) (*Task, error)
}

func (m *MockTaskQueue) AddTask(urlInfo *URLInfo) (*Task, error) {
	if m.AddTaskFunc != nil {
		return m.AddTaskFunc(urlInfo)
	}
	return nil, errors.New("AddTask function not implemented")
}

func (m *MockTaskQueue) StopTask(id int) (*Task, error) {
	if m.StopTaskFunc != nil {
		return m.StopTaskFunc(id)
	}
	return nil, errors.New("StopTask function not implemented")
}

func TestStartComputation(t *testing.T) {
	mockTaskQueue := &MockTaskQueue{
		AddTaskFunc: func(urlInfo *URLInfo) (*Task, error) {
			return &Task{ID: urlInfo.ID, URL: urlInfo.URL}, nil
		},
	}

	mockURLManager := &MockURLManager{
		GetURLInfoFunc: func(id int) *URLInfo {
			return &URLInfo{ID: id, URL: "http://example.com", State: Stopped}
		},
		GetURLStateFunc: func(id int) URLState {
			if id == 1 {
				return Stopped
			}
			return Pending
		},
		UpdateURLStateFunc: func(id int, state URLState) {
			// Mock state update
		},
	}

	logger := logrus.New()
	app := &application{
		taskQueue:  mockTaskQueue,
		urlManager: mockURLManager,
		logger:     logger,
	}

	reqBody := bytes.NewBufferString(`{"id": 1}`)
	req, err := http.NewRequest(http.MethodPost, "/startComputation", reqBody)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.startComputation)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := map[string]interface{}{"id": float64(1), "state": "pending"}
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("error decoding response body: %v", err)
	}

	for k, v := range expected {
		if response[k] != v {
			t.Errorf("handler returned unexpected body: got %v want %v", response, expected)
		}
	}
}
