package main

import (
	"fmt"
	"net/http"
	"os"

	"backend/internal/auth"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

const port = 8080

type application struct {
	authenticator auth.Authenticator
	logger        *logrus.Logger
}

func main() {
	env := os.Getenv("YOURAPP_ENV")
	if env == "" {
		env = "development"
	}

	_ = godotenv.Load(".env." + env + ".local")
	if env != "test" {
		_ = godotenv.Load(".env.local")
	}
	_ = godotenv.Load(".env." + env)
	_ = godotenv.Load() // The Original .env

	// Initialize Authenticator
	authenticator := auth.NewJWTAuthenticator("secret-key")

	// Initialize Logger
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	// Create an instance of the application struct
	app := &application{
		authenticator: authenticator,
		logger:        logger,
	}

	// Test if environment variables are fetched correctly
	fmt.Println("JWT_SECRET:", getEnv("JWT_SECRET", ""))
	// fmt.Println("Domain:", app.Domain)

	logger.Println("Starting application on port", port)

	// start webserver already production ready
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		logger.Fatal(err)
	}

}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
