package auth

import "github.com/golang-jwt/jwt/v5"

type Authenticator interface {
	GenerateAccessToken(claims jwt.Claims) (string, error)
	ValidateAccessToken(token string) (*jwt.Token, error)
	GenerateRefreshToken() (string, error)
}
