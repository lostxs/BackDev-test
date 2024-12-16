package main

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lostxs/BackDev-test/internal/auth"
	"github.com/lostxs/BackDev-test/internal/store"
)

// TODO: Implement logger
type app struct {
	config        config
	store         store.Storage
	authenticator auth.Authenticator
}

type config struct {
	addr string
	db   dbConfig
	auth authConfig
}

type dbConfig struct {
	uri          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

type authConfig struct {
	accessToken accessTokenConfig
}

type accessTokenConfig struct {
	secret string
	exp    time.Duration
}

func (a *app) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Get("/tokens", a.createTokensHandler)
			r.With(a.AccessTokenMiddleware).Get("/refresh", a.refreshTokensHandler)
		})
	})

	return r
}

func (a *app) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         a.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("Starting server on %s", a.config.addr)

	return srv.ListenAndServe()
}
