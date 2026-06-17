package main

import (
	"context"
	"log"

	"spike/01-db-access/ent/ent"
	"spike/01-db-access/scenario"

	"entgo.io/ent/dialect"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://blog:blog@localhost:5432/blog_ent?sslmode=disable"

	client, err := ent.Open(dialect.Postgres, dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	// ent's automatic migration: build the schema from the generated graph.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatal(err)
	}
	if err := scenario.Run(ctx, New(client)); err != nil {
		log.Fatal(err)
	}
}
