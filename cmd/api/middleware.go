package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lostxs/BackDev-test/internal/store"
)

type contextKey string

const (
	userCtx      contextKey = "user"
	ipAddressCtx contextKey = "ip_address"
)

func (a *app) AccessTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			a.unauthorizedException(w, r, fmt.Errorf("authorization header is missing"))
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			a.unauthorizedException(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := a.authenticator.ValidateAccessToken(token)
		if err != nil {
			a.unauthorizedException(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, ok := claims["sub"].(string)
		if !ok {
			a.unauthorizedException(w, r, fmt.Errorf("sub claim is missing"))
			return
		}

		tokenIPAddress, ok := claims["ip_address"].(string)
		if !ok {
			a.unauthorizedException(w, r, fmt.Errorf("ip_address claim is missing"))
			return
		}

		ctx := r.Context()

		user, err := a.getUser(ctx, userID)
		if err != nil {
			a.unauthorizedException(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		ctx = context.WithValue(ctx, ipAddressCtx, tokenIPAddress)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *app) getUser(ctx context.Context, userID string) (*store.User, error) {
	user, err := a.store.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
