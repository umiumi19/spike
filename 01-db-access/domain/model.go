package domain

import "time"

// These structs are deliberately free of any ORM/codegen tags.
// Every implementation must convert its own types into these.
// That conversion friction is itself part of what we are comparing.

type User struct {
	ID        int64
	Name      string
	Email     string
	CreatedAt time.Time
}

type Post struct {
	ID        int64
	Title     string
	Body      string
	AuthorID  int64
	Published bool
	CreatedAt time.Time
}

type Tag struct {
	ID   int64
	Name string
}

type Comment struct {
	ID        int64
	PostID    int64
	AuthorID  int64
	Body      string
	CreatedAt time.Time
}

type PostDetail struct {
	Post     Post
	Author   User
	Tags     []Tag
	Comments []Comment
}
