package main

import (
	"context"
	"database/sql"
	"log"

	"spike/01-db-access/scenario"

	_ "github.com/lib/pq"
)

func main() {
	sdb, err := sql.Open("postgres",
		"postgres://blog:blog@localhost:5432/blog_sqlc?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer sdb.Close()

	// Reminder: schema.sql must already be applied to blog_sqlc (see README).
	if err := scenario.Run(context.Background(), New(sdb)); err != nil {
		log.Fatal(err)
	}
}
