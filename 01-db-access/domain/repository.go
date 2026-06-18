package domain

import "context"

// Repository is the spine of the comparison: gorm, ent, and sqlc each
// implement this exact interface, and cmd/* runs the same scenario against
// each one. If the three implementations produce the same output, they are
// behaving equivalently.
type Repository interface {
	// Minimal insert.
	CreateUser(ctx context.Context, name, email string) (*User, error)

	// Transaction + many-to-many + find-or-create of tags.
	CreatePostWithTags(ctx context.Context, authorID int64, title, body string, published bool, tags []string) (*Post, error)

	// Insert with foreign keys.
	AddComment(ctx context.Context, postID, authorID int64, body string) (*Comment, error)

	// Eager load: Post + Author + Tags + Comments. Watch how each tool emits SQL (N+1?).
	GetPostDetail(ctx context.Context, postID int64) (*PostDetail, error)

	// WHERE + ORDER BY + LIMIT/OFFSET.
	ListPublishedPosts(ctx context.Context, limit, offset int) ([]Post, error)

	// Aggregation.
	CountPostsByAuthor(ctx context.Context, authorID int64) (int, error)
}
