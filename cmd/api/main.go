package main

import (
	"fmt"
	"net/http"

	"strconv"
	"time"

	"backend/internal/auth"
	"backend/internal/services"
	"backend/internal/utils"

	"github.com/sirupsen/logrus"
)

func main() {
	// load environment vars
	utils.LoadEnvFiles()

	jwtSecret := utils.GetEnv("JWT_SECRET", "default-secret")
	if jwtSecret == "" {
		logrus.Fatal("JWT_SECRET environment variable is required but not set")
	}

	allowedOrigin := utils.GetEnv("ALLOWED_ORIGIN", "")
	if allowedOrigin == "" {
		logrus.Fatal("ALLOWED_ORIGIN environment variable is required but not set")
	}

	portStr := utils.GetEnv("PORT", "")
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		logrus.Fatalf("Invalid port number: %v", portStr)
	}

	workersStrt := utils.GetEnv("WORKER_COUNT", "-1")
	workers, err := strconv.Atoi(workersStrt)
	if err != nil || workers < 1 || workers > 100 {
		logrus.Fatalf("Invalid worker count: %v", workers)
	}

	authenticator := auth.NewJWTAuthenticator(jwtSecret)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	urlManager := services.NewURLManager()
	client := &http.Client{Timeout: 10 * time.Second}
	pageAnalyzer := services.NewPageAnalyzer(client, logger)
	taskQueue := services.NewTaskQueue(workers, urlManager, pageAnalyzer, logger)

	app := &application{
		authenticator: authenticator,
		logger:        logger,
		urlManager:    urlManager,
		taskQueue:     taskQueue,
	}

	logger.Println("Starting application on port", port)
	logger.Infof("Workers count: %d", workers)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: app.routes(),
	}

	go func() {
		// for production => ListenAndServeTLS
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Could not listen on %s: %v\n", srv.Addr, err)
		}
	}()

	utils.GracefulShutdown(srv, logger)

}
