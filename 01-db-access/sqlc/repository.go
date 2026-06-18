package main

import (
	"context"
	"database/sql"
	"spike/01-db-access/domain"
	"spike/01-db-access/sqlc/gen"
)

type Repo struct {
	db *sql.DB
	q *gen.Queries
}

func New(sqlDB *sql.DB) *Repo {
	return &Repo{db: sqlDB, q: gen.New(sqlDB)}
}

func (r *Repo) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	u, err := r.q.CreateUser(ctx, gen.CreateUserParams{Name: name, Email: email})
	if err != nil {
		return nil, err
	}
	return toDomainUser(u), nil
}

func (r *Repo) CreatePostWithTags(ctx context.Context, authorID int64, title, body string, published bool, tagNames []string) (*domain.Post, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := r.q.WithTx(tx)

	p, err := qtx.CreatePost(ctx, gen.CreatePostParams{
		Title: title, Body: body, AuthorID:  authorID, Published:  published,
	})
	if err != nil {
		return nil, err
	}

	for _, name := range tagNames {
		t, err := qtx.GetOrCreateTag(ctx, name)
		if err != nil {
			return nil, err
		}
		if err := qtx.AttachTag(ctx, gen.AttachTagParams{
			PostID: p.ID, TagID: t.ID,
		}); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return toDomainPost(p), nil
}

func (r *Repo) AddComment(ctx context.Context, postID, authorID int64, body string) (*domain.Comment, error) {
	c, err := r.q.AddComment(ctx, gen.AddCommentParams{
		PostID: postID, AuthorID: authorID, Body: body,
	})
	if err != nil {
		return nil, err
	}

	return toDomainComment(c), nil
}

func (r *Repo) GetPostDetail(ctx context.Context, postID int64) (*domain.PostDetail, error) {
	p, err := r.q.GetPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	author, err := r.q.GetUser(ctx, p.AuthorID)
	if err != nil {
		return nil, err
	}
	tags, err := r.q.ListTagsByPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	comments, err := r.q.ListCommentByPost(ctx, postID)
	if err != nil {
		return nil, err
	}
	return toDomainPostDetail(p, author, tags, comments), nil
}

func (r *Repo) ListPublishedPosts(ctx context.Context, limit, offset int) ([]domain.Post, error) {
	ps, err := r.q.ListPublishedPosts(ctx, gen.ListPublishedPostsParams{
		Limit: int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	out := make([]domain.Post, 0, len(ps))
	for _, p := range ps {
		out = append(out, *toDomainPost(p))
	}

	return out, nil
}

func (r *Repo) CountPostsByAuthor(ctx context.Context, authorID int64) (int, error) {
	n, err := r.q.CountPostsByAuthor(ctx, authorID)
	if err != nil {
		return 0, err
	}
	return int(n), nil
}

func toDomainUser(u gen.User) *domain.User {
	return &domain.User{ID: u.ID, Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt}
}

func toDomainPost(p gen.Post) *domain.Post {
	return &domain.Post{
		ID: p.ID,
		Title: p.Title,
		Body: p.Body,
		AuthorID: p.AuthorID,
		Published: p.Published,
		CreatedAt: p.CreatedAt,
	}
}

func toDomainComment(c gen.Comment) *domain.Comment {
	return &domain.Comment{
		ID: c.ID,
		PostID: c.PostID,
		AuthorID: c.AuthorID,
		Body: c.Body,
		CreatedAt: c.CreatedAt,
	}
}

func toDomainPostDetail(post gen.Post, author gen.User, tags []gen.Tag, comments []gen.Comment) *domain.PostDetail {
	detail := &domain.PostDetail{
		Post: *toDomainPost(post),
		Author: *toDomainUser(author),
	}

	for _, t := range tags {
		detail.Tags = append(detail.Tags, domain.Tag{ID: t.ID, Name: t.Name})
	}

	for _, c := range comments {
		detail.Comments = append(detail.Comments, *toDomainComment(c))
	}
	return detail
}