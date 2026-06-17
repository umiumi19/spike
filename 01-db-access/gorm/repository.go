package main

import (
	"context"
	"time"

	"spike/01-db-access/domain"
	"spike/01-db-access/gorm/model"
	"spike/01-db-access/gorm/query"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
	q  *query.Query
}

func New(db *gorm.DB) *Repo {
	return &Repo{db: db, q: query.Use(db)}
}

func (r *Repo) Migrate() error {
	return r.db.AutoMigrate(&model.User{}, &model.Post{}, &model.Tag{}, &model.Comment{})
}

func (r *Repo) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	u := &model.User{Name: name, Email: email, CreatedAt: time.Now()}
	if err := r.q.User.WithContext(ctx).Create(u); err != nil {
		return nil, err
	}

	return toDomainUser(u), nil
}

func (r *Repo) CreatePostWithTags(ctx context.Context, authorID int64, title, body string, published bool, tagNames []string) (*domain.Post, error) {
	post := &model.Post{Title: title, Body: body, AuthorID: authorID, Published: published, CreatedAt: time.Now()}

	err := r.q.Transaction(func(tx *query.Query) error {
		if err := tx.Post.WithContext(ctx).Create(post); err != nil {
			return err
		}

		for _, name := range tagNames {
			tag, err := tx.Tag.WithContext(ctx).
				Where(tx.Tag.Name.Eq(name)).
				FirstOrCreate()
			if err != nil {
				return err
			}
			if err := tx.Post.Tags.Model(post).Append(tag); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return toDomainPost(post), nil
}

func (r *Repo) AddComment(ctx context.Context, postID, authorID int64, body string) (*domain.Comment, error) {
	c := &model.Comment{PostID: postID, AuthorID: authorID, Body: body, CreatedAt: time.Now()}
	if err := r.q.Comment.WithContext(ctx).Create(c); err != nil {
		return nil, err
	}

	return toDomainComment(c), nil

}

func (r *Repo) GetPostDetail(ctx context.Context, postID int64) (*domain.PostDetail, error) {
	post, err := r.q.Post.WithContext(ctx).
		Where(r.q.Post.ID.Eq(postID)).
		Preload(r.q.Post.Author).
		Preload(r.q.Post.Tags).
		Preload(r.q.Post.Comments).
		First()
	if err != nil {
		return nil, err
	}
	return toDomainPostDetail(post), nil
}

func (r *Repo) ListPublishedPosts(ctx context.Context, limit, offset int) ([]domain.Post, error) {
	posts, err := r.q.Post.WithContext(ctx).
		Where(r.q.Post.Published.Is(true)).
		Order(r.q.Post.CreatedAt.Desc()).
		Limit(limit).
		Offset(offset).
		Find()
	if err != nil {
		return nil, err
	}

	out := make([]domain.Post, 0, len(posts))
	for _, p := range posts {
		out = append(out, *toDomainPost(p))
	}

	return out, nil
}

func (r *Repo) CountPostsByAuthor(ctx context.Context, authorID int64) (int, error) {
	n, err := r.q.Post.WithContext(ctx).
		Where(r.q.Post.AuthorID.Eq(authorID)).
		Count()
	return int(n), err
}

func toDomainUser(u *model.User) *domain.User {
	return &domain.User{ID: u.ID, Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt}
}

func toDomainPost(p *model.Post) *domain.Post {
	return &domain.Post{ID: p.ID, Title: p.Title, Body: p.Body, AuthorID: p.AuthorID, Published: p.Published, CreatedAt: p.CreatedAt}
}

func toDomainComment(c *model.Comment) *domain.Comment {
	return &domain.Comment{ID: c.ID, PostID: c.PostID, AuthorID: c.AuthorID, Body: c.Body, CreatedAt: c.CreatedAt}
}

func toDomainPostDetail(p *model.Post) *domain.PostDetail {
	detail := &domain.PostDetail{Post: *toDomainPost(p), Author: *toDomainUser(&p.Author)}

	for _, t := range p.Tags {
		detail.Tags = append(detail.Tags, domain.Tag{ID: t.ID, Name: t.Name})
	}
	for _, c := range p.Comments {
		detail.Comments = append(detail.Comments, *toDomainComment(&c))
	}

	return detail
}

