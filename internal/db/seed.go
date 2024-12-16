package db

import (
	"context"
	"database/sql"
	"log"

	"github.com/lostxs/BackDev-test/internal/store"
)

var emails = []string{
	"admin@example.com",
	"user@example.com",
	"guest@example.com",
}

func Seed(store store.Storage, db *sql.DB) {
	ctx := context.Background()

	users := generateUsers(len(emails))
	tx, _ := db.BeginTx(ctx, nil)

	for _, user := range users {
		if err := store.Users.Create(ctx, tx, user); err != nil {
			_ = tx.Rollback()
			log.Println("Error creating user:", err)
			return
		}
	}

	tx.Commit()

	log.Println("Seeding complete")
}

func generateUsers(num int) []*store.User {
	users := make([]*store.User, num)

	for i := 0; i < num; i++ {
		users[i] = &store.User{
			Email: emails[i%len(emails)],
		}
	}

	return users
}
