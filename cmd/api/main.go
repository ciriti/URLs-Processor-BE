package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"backend/internal/auth"

	"github.com/joho/godotenv"
)

const port = 8080

type application struct {
    authenticator auth.Authenticator
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

	// Create an instance of the application struct
	app := &application{
		authenticator: authenticator,
	}

	// Test if environment variables are fetched correctly
	fmt.Println("JWT_SECRET:", getEnv("JWT_SECRET", ""))
	// fmt.Println("Domain:", app.Domain)

	log.Println("Starting application on port", port)

	// start webserver already production ready
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}

}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
