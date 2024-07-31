package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const port = 8080

type application struct {
	DSN    string
	Domain string
	// DB           repository.DatabaseRepo
	// auth         Authenticator // Use the interface here
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string
	APIKey       string
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

	// Create an instance of the application struct
	app := application{
		DSN:          getEnv("DSN", ""),
		Domain:       getEnv("DOMAIN", ""),
		JWTSecret:    getEnv("JWT_SECRET", ""),
		JWTIssuer:    getEnv("JWT_ISSUER", ""),
		JWTAudience:  getEnv("JWT_AUDIENCE", ""),
		CookieDomain: getEnv("COOKIE_DOMAIN", ""),
		APIKey:       getEnv("API_KEY", ""),
	}

	// Test if environment variables are fetched correctly
	fmt.Println("DSN:", app.DSN)
	fmt.Println("Domain:", app.Domain)
	fmt.Println("JWTSecret:", app.JWTSecret)
	fmt.Println("JWTIssuer:", app.JWTIssuer)
	fmt.Println("JWTAudience:", app.JWTAudience)
	fmt.Println("CookieDomain:", app.CookieDomain)
	fmt.Println("APIKey:", app.APIKey)

}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
