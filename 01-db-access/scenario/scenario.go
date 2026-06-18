package scenario

import (
	"context"
	"fmt"
	"time"

	"spike/01-db-access/domain"
)

// Run executes the same story against any Repository implementation.
// Email is timestamped so re-runs don't trip the UNIQUE(email) constraint.
func Run(ctx context.Context, repo domain.Repository) error {
	email := fmt.Sprintf("alice+%d@example.com", time.Now().UnixNano())

	alice, err := repo.CreateUser(ctx, "Alice", email)
	if err != nil {
		return fmt.Errorf("CreateUser: %w", err)
	}
	fmt.Printf("user:   id=%d name=%s\n", alice.ID, alice.Name)

	post, err := repo.CreatePostWithTags(ctx, alice.ID, "Hello", "world body", true, []string{"go", "db"})
	if err != nil {
		return fmt.Errorf("CreatePostWithTags: %w", err)
	}
	fmt.Printf("post:   id=%d title=%s published=%v\n", post.ID, post.Title, post.Published)

	if _, err := repo.AddComment(ctx, post.ID, alice.ID, "nice post"); err != nil {
		return fmt.Errorf("AddComment: %w", err)
	}

	detail, err := repo.GetPostDetail(ctx, post.ID)
	if err != nil {
		return fmt.Errorf("GetPostDetail: %w", err)
	}
	fmt.Printf("detail: title=%q author=%s tags=%d comments=%d\n",
		detail.Post.Title, detail.Author.Name, len(detail.Tags), len(detail.Comments))

	published, err := repo.ListPublishedPosts(ctx, 10, 0)
	if err != nil {
		return fmt.Errorf("ListPublishedPosts: %w", err)
	}
	fmt.Printf("list:   published posts (max 10) = %d\n", len(published))

	count, err := repo.CountPostsByAuthor(ctx, alice.ID)
	if err != nil {
		return fmt.Errorf("CountPostsByAuthor: %w", err)
	}
	fmt.Printf("count:  posts by alice = %d\n", count)

	fmt.Println("OK")
	return nil
}
