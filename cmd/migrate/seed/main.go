package main

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/lostxs/BackDev-test/internal/db"
	"github.com/lostxs/BackDev-test/internal/env"
	"github.com/lostxs/BackDev-test/internal/store"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	uri := env.GetString("DATABASE_URI", "postgres://postgres:postgres@localhost:5432/backdev?sslmode=disable")
	conn, err := db.New(uri, 10, 10, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewPostgresStorage(conn)

	db.Seed(store, conn)
}
