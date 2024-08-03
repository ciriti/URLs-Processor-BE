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

type MockURLManager struct {
	AddURLFunc              func(url string) *URLInfo
	UpdateURLStateFunc      func(id int, state URLState)
	UpdateProcessedDataFunc func(id int, data *DataInfo)
	GetURLInfoFunc          func(id int) *URLInfo
	GetAllURLsFunc          func() []*URLInfo
	NextIDFunc              func() int
	GetURLStateFunc         func(id int) URLState
}

func (m *MockURLManager) AddURL(url string) *URLInfo {
	if m.AddURLFunc != nil {
		return m.AddURLFunc(url)
	}
	return nil
}

func (m *MockURLManager) UpdateURLState(id int, state URLState) {
	if m.UpdateURLStateFunc != nil {
		m.UpdateURLStateFunc(id, state)
	}
}

func (m *MockURLManager) UpdateProcessedData(id int, data *DataInfo) {
	if m.UpdateProcessedDataFunc != nil {
		m.UpdateProcessedDataFunc(id, data)
	}
}

func (m *MockURLManager) GetURLInfo(id int) *URLInfo {
	if m.GetURLInfoFunc != nil {
		return m.GetURLInfoFunc(id)
	}
	return nil
}

func (m *MockURLManager) GetAllURLs() []*URLInfo {
	if m.GetAllURLsFunc != nil {
		return m.GetAllURLsFunc()
	}
	return nil
}

func (m *MockURLManager) nextID() int {
	if m.NextIDFunc != nil {
		return m.NextIDFunc()
	}
	return 0
}

func (m *MockURLManager) GetURLState(id int) URLState {
	if m.GetURLStateFunc != nil {
		return m.GetURLStateFunc(id)
	}
	return ""
}

type MockTaskQueue struct {
	AddTaskFunc  func(urlInfo *URLInfo) (*Task, error)
	StopTaskFunc func(id int) (*Task, error)
	GetTaskFunc  func(id int) (*Task, error)
	ContainsFunc func(id int) bool
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

func (m *MockTaskQueue) GetTask(id int) (*Task, error) {
	if m.GetTaskFunc != nil {
		return m.GetTaskFunc(id)
	}
	return nil, errors.New("GetTask function not implemented")
}

func (m *MockTaskQueue) Contains(id int) bool {
	if m.ContainsFunc != nil {
		return m.ContainsFunc(id)
	}
	return false
}

func TestStartComputation(t *testing.T) {
	mockTaskQueue := &MockTaskQueue{
		AddTaskFunc: func(urlInfo *URLInfo) (*Task, error) {
			return &Task{ID: urlInfo.ID, URL: urlInfo.URL}, nil
		},
		GetTaskFunc: func(id int) (*Task, error) {
			if id == 1 {
				return &Task{ID: 1, URL: "http://example.com"}, nil
			}
			return nil, errors.New("task not found")
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
