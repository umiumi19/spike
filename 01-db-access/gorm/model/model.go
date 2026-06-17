package model

import "time"

type User struct {
	ID        int64  `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"uniqueIndex"`
	CreatedAt time.Time
	Posts     []Post `gorm:"foreignKey:AuthorID"`
}

type Post struct {
	ID        int64 `gorm:"primaryKey"`
	Title     string
	Body      string
	AuthorID  int64
	Author    User
	Published bool
	CreatedAt time.Time
	Tags      []Tag     `gorm:"many2many:post_tags;"`
	Comments  []Comment
}

type Tag struct {
	ID    int64  `gorm:"primaryKey"`
	Name  string `gorm:"uniqueIndex"`
	Posts []Post `gorm:"many2many:post_tags;"`
}

type Comment struct {
	ID        int64 `gorm:"primaryKey"`
	PostID    int64
	AuthorID  int64
	Body      string
	CreatedAt time.Time
}
