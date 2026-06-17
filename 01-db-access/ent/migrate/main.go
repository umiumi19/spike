//go:build ignore

package main

import (
	"context"
	"log"
	"os"

	"ariga.io/atlas/sql/sqltool"
	"entgo.io/ent/dialect/sql/schema"
	_ "github.com/lib/pq"

	entmigrate "spike/01-db-access/ent/ent/migrate"
)

const devURL = "postgres://blog:blog@localhost:5432/blog_ent?sslmode=disable"

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: go run migrate/main.go <migration-name>")
	}

	dir, err := sqltool.NewGolangMigrateDir("ent/migrations")
	if err != nil {
		log.Fatal("failed to open migration dir:", err)
	}

	opts := []schema.MigrateOption{
		schema.WithDir(dir),
		schema.WithMigrationMode(schema.ModeReplay),
		schema.WithDropColumn(true),
		schema.WithDropIndex(true),
	}

	err = entmigrate.NamedDiff(context.Background(), devURL, os.Args[1], opts...)
	if err != nil {
		log.Fatal("failed to generate migration:", err)
	}
}
