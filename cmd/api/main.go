package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"backend/internal/auth"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type application struct {
	authenticator auth.Authenticator
	logger        *logrus.Logger
	urlManager    URLManagerInterface
	taskQueue     TaskQueueInterface
}

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	_ = godotenv.Load(".env." + env + ".local")
	if env != "test" {
		_ = godotenv.Load(".env.local")
	}
	_ = godotenv.Load(".env." + env)
	_ = godotenv.Load()

	jwtSecret := getEnv("JWT_SECRET", "default-secret")
	if jwtSecret == "" {
		logrus.Fatal("JWT_SECRET environment variable is required but not set")
	}

	allowedOrigin := getEnv("ALLOWED_ORIGIN", "")
	if allowedOrigin == "" {
		logrus.Fatal("ALLOWED_ORIGIN environment variable is required but not set")
	}

	portStr := getEnv("PORT", "")
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 1 || port > 65535 {
		logrus.Fatalf("Invalid port number: %v", portStr)
	}

	authenticator := auth.NewJWTAuthenticator(jwtSecret)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	urlManager := NewURLManager()

	workersStrt := getEnv("WORKER_COUNT", "-1")
	workers, err := strconv.Atoi(workersStrt)
	logger.Infof("Workers count: %d", workers)
	if err != nil || workers < 1 || workers > 100 {
		logrus.Fatalf("Invalid worker count: %v", workers)
	}
	taskQueue := NewTaskQueue(workers, urlManager, logger)

	app := &application{
		authenticator: authenticator,
		logger:        logger,
		urlManager:    urlManager,
		taskQueue:     taskQueue,
	}

	logger.Println("Starting application on port", port)

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

	gracefulShutdown(srv, logger)

}

// all in-flight requests are completed before the server stops
func gracefulShutdown(srv *http.Server, logger *logrus.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Println("Server exiting")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
