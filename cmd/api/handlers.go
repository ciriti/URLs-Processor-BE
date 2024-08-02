package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"
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

	var payload struct {
		User string `json:"user"`
		Pass string `json:"pass"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.logger.WithError(err).Error("error decoding JSON request body")
		err = app.errorJSON(w, errors.New("invalid request payload"), http.StatusBadRequest)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
		return
	}

	app.logger.Infof("Received authentication request - user: %s, pass: %s", payload.User, payload.Pass)

	if payload.User == "" || payload.Pass == "" {
		err := app.errorJSON(w, errors.New("username and password are required"), http.StatusBadRequest)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
		return
	}

	if !app.authenticator.ValidateUserCredentials(payload.User, payload.Pass) {
		err := app.errorJSON(w, errors.New(http.StatusText(http.StatusUnauthorized)), http.StatusUnauthorized)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
		return
	}

	tokenString, err := app.authenticator.GenerateToken(payload.User)
	if err != nil {
		app.logger.WithError(err).Error("error generating token")
		err := app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "jwtToken",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	if err := app.writeJSON(w, http.StatusOK, map[string]string{"token": tokenString}); err != nil {
		app.logger.WithError(err).Error("error writing JSON response")
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "jwtToken",
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	if err := app.writeJSON(w, http.StatusOK, map[string]string{"message": "logout successful"}); err != nil {
		app.logger.WithError(err).Error("error writing JSON response")
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
	}
}

func (app *application) adminEndpoint(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("This is a protected endpoint"))
	if err != nil {
		app.logger.WithError(err).Error("error writing JSON response")
	}
}

func (app *application) regularEndpoint(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("This is a regular endpoint"))
	if err != nil {
		app.logger.WithError(err).Error("error writing JSON response")
	}
}

func (app *application) addURLs(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		URLs []string `json:"urls"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.logger.WithError(err).Error("error decoding JSON request body")
		err = app.errorJSON(w, errors.New("invalid request payload"), http.StatusBadRequest)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
		return
	}

	for _, url := range payload.URLs {
		app.urlManager.AddURL(url)
		app.logger.Infof("Adding URL: %s", url)
	}

	if err := app.writeJSON(w, http.StatusOK, map[string]string{"message": "URLs added successfully"}); err != nil {
		app.logger.WithError(err).Error("error writing JSON response")
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
	}
}

func (app *application) getURL(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		err := app.errorJSON(w, errors.New("missing url id parameter"), http.StatusBadRequest)
		if err != nil {
			app.logger.Println("error writing JSON response:", err)
		}
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		err := app.errorJSON(w, errors.New("invalid url id parameter"), http.StatusBadRequest)
		if err != nil {
			app.logger.Println("error writing JSON response:", err)
		}
		return
	}

	urlInfo := app.urlManager.GetURLInfo(id)
	if urlInfo == nil {
		err := app.errorJSON(w, errors.New("URL not found"), http.StatusNotFound)
		if err != nil {
			app.logger.Println("error writing JSON response:", err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, urlInfo); err != nil {
		app.logger.Println("error writing JSON response:", err)
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			app.logger.Println("error writing JSON response:", err)
		}
	}
}

func (app *application) getAllURLs(w http.ResponseWriter, r *http.Request) {
	urls := app.urlManager.GetAllURLs()

	if err := app.writeJSON(w, http.StatusOK, urls); err != nil {
		app.logger.WithError(err).Error("error writing JSON response")
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
	}
}

func (app *application) startComputation(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		ID int `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.logger.WithError(err).Error("error decoding JSON request body")
		err = app.errorJSON(w, errors.New("invalid request payload"), http.StatusBadRequest)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
		return
	}

	urlInfo := app.urlManager.GetURLInfo(payload.ID)
	if urlInfo == nil {
		app.logger.WithError(err).Error("URL not found")
		err = app.errorJSON(w, errors.New("URL not found"), http.StatusNotFound)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
		return
	}

	task, err := app.taskQueue.AddTask(urlInfo)
	if err != nil {
		app.logger.WithError(err).Error("task already in progress")
		err = app.errorJSON(w, errors.New("task already in progress"), http.StatusConflict)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, map[string]interface{}{"task_id": task.ID, "state": urlInfo.State}); err != nil {
		app.logger.WithError(err).Error("error writing JSON response")
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
	}
}

func (app *application) checkStatus(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		err := app.errorJSON(w, errors.New("missing url id parameter"), http.StatusBadRequest)
		if err != nil {
			app.logger.Println("error writing JSON response:", err)
		}
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		err := app.errorJSON(w, errors.New("invalid url id parameter"), http.StatusBadRequest)
		if err != nil {
			app.logger.Println("error writing JSON response:", err)
		}
		return
	}

	urlInfo := app.urlManager.GetURLInfo(id)
	if urlInfo == nil {
		err := app.errorJSON(w, errors.New("URL not found"), http.StatusNotFound)
		if err != nil {
			app.logger.Println("error writing JSON response:", err)
		}
		return
	}

	app.logger.Infof("Task ID: %d, Status: %s", urlInfo.ID, urlInfo.State)

	if err := app.writeJSON(w, http.StatusOK, urlInfo); err != nil {
		app.logger.Println("error writing JSON response:", err)
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			app.logger.Println("error writing JSON response:", err)
		}
	}
}

func (app *application) stopComputation(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		TaskID int `json:"task_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.logger.WithError(err).Error("error decoding JSON request body")
		err = app.errorJSON(w, errors.New("invalid request payload"), http.StatusBadRequest)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
		return
	}

	app.logger.Infof("Stopping task - task_id: %d", payload.TaskID)
	app.taskQueue.StopTask(payload.TaskID)

	if err := app.writeJSON(w, http.StatusOK, map[string]string{"message": "Task stop signal sent successfully"}); err != nil {
		app.logger.WithError(err).Error("error writing JSON response")
		err = app.errorJSON(w, err, http.StatusInternalServerError)
		if err != nil {
			app.logger.WithError(err).Error("error writing JSON response")
		}
	}
}
