package auth

import (
	"github.com/golang-jwt/jwt/v5"
)

type TestAuthenticator struct{}

const secret = "test"

func NewTestAuthenticator() *TestAuthenticator {
	return &TestAuthenticator{}
}

func (a *TestAuthenticator) GenerateAccessToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString, nil
}

func (a *TestAuthenticator) ValidateAccessToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
}

func (a *TestAuthenticator) GenerateRefreshToken() (string, error) {
	return "test", nil
}
