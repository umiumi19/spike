package main

import (
	"context"
	"log"

	"spike/01-db-access/scenario"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "postgres://blog:blog@localhost:5432/blog_gorm?sslmode=disable"

	gdb, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Show every SQL statement. Great for spotting Preload's N+1 behavior.
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal(err)
	}

	repo := New(gdb)
	if err := repo.Migrate(); err != nil {
		log.Fatal(err)
	}
	if err := scenario.Run(context.Background(), repo); err != nil {
		log.Fatal(err)
	}
}
