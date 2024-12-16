package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lostxs/BackDev-test/internal/auth"
	"github.com/lostxs/BackDev-test/internal/store"
)

func newTestApplication(t *testing.T, cfg config) *app {
	t.Helper()

	mockStore := store.NewMockStore()

	testAuth := auth.NewTestAuthenticator()

	return &app{
		config:        cfg,
		store:         mockStore,
		authenticator: testAuth,
	}
}

func executeRequest(req *http.Request, mux http.Handler) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d", expected, actual)
	}
}
