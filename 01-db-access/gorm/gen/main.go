package main

import (
	"spike/01-db-access/gorm/model"

	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "../query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery,
	})

	g.ApplyBasic(
		model.User{},
		model.Post{},
		model.Tag{},
		model.Comment{},
	)

	g.Execute()
}

