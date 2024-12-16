package main

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lostxs/BackDev-test/internal/auth"
	"github.com/lostxs/BackDev-test/internal/store"
)

// TODO: Refractor code

func TestCreateTokensHandler(t *testing.T) {
	cfg := config{}

	app := newTestApplication(t, cfg)

	mockUserStore := app.store.Users.(*store.MockUserStore)
	mockUserStore.Create(context.Background(), nil, &store.User{
		ID:    "86990727-379a-42ea-a71d-69179969e777",
		Email: "test@test.com",
	})

	mockSessionStore := app.store.Sessions.(*store.MockSessionStore)
	mockSessionStore.Upsert(context.Background(), &store.Session{
		ID:     "ce2c7489-837a-4910-84b8-cff4e70248a5",
		UserID: "86990727-379a-42ea-a71d-69179969e777",
	})

	mux := app.mount()

	t.Run("should return 400 if user_id is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/auth/tokens", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusBadRequest, rr.Code)

		expected := "user_id not provided"
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("expected error message %q, got %q", expected, rr.Body.String())
		}
	})

	t.Run("should return 400 if user_id is invalid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/auth/tokens?user_id=invalid-user-id", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusBadRequest, rr.Code)

		if !strings.Contains(rr.Body.String(), store.ErrInvalidUserID.Error()) {
			t.Errorf("expected error message %q, got %q", store.ErrInvalidUserID.Error(), rr.Body.String())
		}
	})

	t.Run("should return 401 if user is not found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/auth/tokens?user_id=1e2e06f9-a42f-4e9e-a5e0-f2f376e70dc6", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)

		if !strings.Contains(rr.Body.String(), store.ErrUserNotFound.Error()) {
			t.Errorf("expected error message %q, got %q", store.ErrUserNotFound.Error(), rr.Body.String())
		}
	})

	t.Run("should generate access token", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/auth/tokens?user_id=86990727-379a-42ea-a71d-69179969e777", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.RemoteAddr = "127.0.0.1:8080"

		app.authenticator.GenerateAccessToken(jwt.MapClaims{
			"sub":        "86990727-379a-42ea-a71d-69179969e777",
			"ip_address": req.RemoteAddr,
		})
	})

	t.Run("should create session and refresh token", func(t *testing.T) {
		testRefreshToken, err := app.authenticator.GenerateRefreshToken()
		if err != nil {
			t.Fatal(err)
		}

		hash, err := hashValue(testRefreshToken)
		if err != nil {
			t.Fatal(err)
		}

		mockSessionStore.Upsert(context.Background(), &store.Session{
			ID:               "ce2c7489-837a-4910-84b8-cff4e70248a5",
			UserID:           "86990727-379a-42ea-a71d-69179969e777",
			RefreshTokenHash: string(hash),
		})
	})

	t.Run("should set refresh token in cookie", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/auth/tokens?user_id=86990727-379a-42ea-a71d-69179969e777", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)

		cookies := rr.Result().Cookies()
		var refreshTokenCookie *http.Cookie
		for _, c := range cookies {
			if c.Name == "refresh_token" {
				refreshTokenCookie = c
				break
			}
		}

		if refreshTokenCookie == nil {
			t.Fatalf("expected refresh_token cookie to be set")
		}
		if refreshTokenCookie.Path != "/" {
			t.Errorf("expected refresh_token cookie Path to be '/', got %q", refreshTokenCookie.Path)
		}
		if !refreshTokenCookie.HttpOnly {
			t.Errorf("expected refresh_token cookie HttpOnly to be true")
		}
	})

	t.Run("should return valid JSON response", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/auth/tokens?user_id=86990727-379a-42ea-a71d-69179969e777", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)

		expected := `"access_token":`
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("expected JSON response to contain %q, got %q", expected, rr.Body.String())
		}

		expectedUserField := `"email":"test@test.com"`
		if !strings.Contains(rr.Body.String(), expectedUserField) {
			t.Errorf("expected JSON response to contain %q, got %q", expectedUserField, rr.Body.String())
		}
	})

	t.Run("should return 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/auth/tokens?user_id=86990727-379a-42ea-a71d-69179969e777", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusOK, rr.Code)
	})
}

func TestRefreshHandler(t *testing.T) {
	cfg := config{}

	app := newTestApplication(t, cfg)

	mockUserStore := app.store.Users.(*store.MockUserStore)
	mockUserStore.Create(context.Background(), nil, &store.User{
		ID:    "86990727-379a-42ea-a71d-69179969e777",
		Email: "test@test.com",
	})

	mockSessionStore := app.store.Sessions.(*store.MockSessionStore)
	mockSessionStore.Upsert(context.Background(), &store.Session{
		ID:               "ce2c7489-837a-4910-84b8-cff4e70248a5",
		UserID:           "86990727-379a-42ea-a71d-69179969e777",
		RefreshTokenHash: hashValueOrFail("valid-refresh-token"),
	})

	testAuthenticator := app.authenticator.(*auth.TestAuthenticator)
	testClaims := jwt.MapClaims{
		"sub":        "86990727-379a-42ea-a71d-69179969e777",
		"ip_address": "127.0.0.1:8080",
		"exp":        time.Now().Add(time.Hour).Unix(),
	}

	testAccessToken, err := testAuthenticator.GenerateAccessToken(testClaims)
	if err != nil {
		t.Fatal(err)
	}

	mux := app.mount()

	t.Run("should return 401 authorization header is missing", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/auth/refresh", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return 401 if access token is invalid", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/auth/refresh", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer invalid-access-token")

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("should return 401 if refresh token is not provided", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/auth/refresh", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testAccessToken)

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)

		expected := "refresh token not provided"
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("expected error message %q, got %q", expected, rr.Body.String())
		}
	})

	t.Run("should return 401 if session not found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/api/auth/refresh", nil)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Authorization", "Bearer "+testAccessToken)
		req.AddCookie(&http.Cookie{
			Name:  "refresh_token",
			Value: "refresh-token",
		})

		mockSessionStore.Delete(context.Background(), "ce2c7489-837a-4910-84b8-cff4e70248a5")

		rr := executeRequest(req, mux)
		checkResponseCode(t, http.StatusUnauthorized, rr.Code)

		expected := "session not found"
		if !strings.Contains(rr.Body.String(), expected) {
			t.Errorf("expected error message %q, got %q", expected, rr.Body.String())
		}
	})

}

func hashValueOrFail(value string) string {
	hash, err := hashValue(value)
	if err != nil {
		panic(err)
	}
	return string(hash)
}
