package auth

import (
	"time"

	"github.com/go-chi/jwtauth"
	"github.com/golang-jwt/jwt/v4"
)

type Authenticator interface {
	ValidateUserCredentials(user, pass string) bool
	GenerateToken(user string) (string, error)
}

type JWTAuthenticator struct {
	TokenAuth *jwtauth.JWTAuth
}

func NewJWTAuthenticator(secret string) *JWTAuthenticator {
	return &JWTAuthenticator{
		TokenAuth: jwtauth.New("HS256", []byte(secret), nil),
	}
}

func (a *JWTAuthenticator) ValidateUserCredentials(user, pass string) bool {
	return user == "admin" && pass == "password"
}

func (a *JWTAuthenticator) GenerateToken(user string) (string, error) {
	_, tokenString, err := a.TokenAuth.Encode(jwt.MapClaims{
		"user": user,
		"exp":  time.Now().Add(time.Hour * 72).Unix(), // 72 hours expiration
	})
	return tokenString, err
}
