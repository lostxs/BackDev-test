package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/lostxs/BackDev-test/internal/auth"
	"github.com/lostxs/BackDev-test/internal/db"
	"github.com/lostxs/BackDev-test/internal/env"
	"github.com/lostxs/BackDev-test/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config{
		addr: env.GetString("APP_ADDR", ":8080"),
		db: dbConfig{
			uri:          env.GetString("DATABASE_URI", "postgres://postgres:postgres@localhost:5432/backdev?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 10),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME", 15*time.Minute),
		},
		auth: authConfig{
			accessToken: accessTokenConfig{
				secret: env.GetString("ACCESS_TOKEN_SECRET", "access_secret"),
				exp:    env.GetDuration("ACCESS_TOKEN_EXP", 15*time.Minute),
			},
		},
	}

	db, err := db.New(
		cfg.db.uri,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("Database connection pool established")

	store := store.NewPostgresStorage(db)

	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.accessToken.secret,
	)

	app := app{
		config:        cfg,
		store:         store,
		authenticator: jwtAuthenticator,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
