package main

import (
	"context"
	"fmt"

	"spike/01-db-access/domain"
	"spike/01-db-access/ent/ent"
	"spike/01-db-access/ent/ent/comment"
	"spike/01-db-access/ent/ent/post"
	"spike/01-db-access/ent/ent/tag"
	"spike/01-db-access/ent/ent/user"
)

type Repo struct {
	client *ent.Client
}

func New(client *ent.Client) *Repo {
	return &Repo{client: client}
}

func (r *Repo) CreateUser(ctx context.Context, name, email string) (*domain.User, error) {
	u, err := r.client.User.Create().
		SetName(name).
		SetEmail(email).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return toDomainUser(u), nil
}

func (r *Repo) CreatePostWithTags(ctx context.Context, authorID int64, title, body string, published bool, tagNames []string) (*domain.Post, error) {
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, err
	}

	tagIDs := make([]int, 0, len(tagNames))
	for _, name := range tagNames {
		t, err := tx.Tag.Query().Where(tag.NameEQ(name)).Only(ctx)
		if ent.IsNotFound(err) {
			t, err = tx.Tag.Create().SetName(name).Save(ctx)
		}
		if err != nil {
			return nil, rollback(tx, err)
		}
		tagIDs = append(tagIDs, t.ID)
	}

	p, err := tx.Post.Create().
		SetTitle(title).
		SetBody(body).
		SetPublished(published).
		SetAuthorID(int(authorID)).
		AddTagIDs(tagIDs...).
		Save(ctx)
	if err != nil {
		return nil, rollback(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	p2, err := r.client.Post.Query().
		Where(post.IDEQ(p.ID)).
		WithAuthor().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return toDomainPost(p2), nil
}

func (r *Repo) AddComment(ctx context.Context, postID, authorID int64, body string) (*domain.Comment, error) {
	c, err := r.client.Comment.Create().
		SetBody(body).
		SetPostID(int(postID)).
		SetAuthorID(int(authorID)).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	c2, err := r.client.Comment.Query().
		Where(comment.IDEQ(c.ID)).
		WithPost().
		WithAuthor().
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return toDomainComment(c2), nil
}

func (r *Repo) GetPostDetail(ctx context.Context, postID int64) (*domain.PostDetail, error) {
	p, err := r.client.Post.Query().
		Where(post.IDEQ(int(postID))).
		WithAuthor().
		WithTags().
		WithComments(func(q *ent.CommentQuery) {
			q.WithAuthor().WithPost()
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	return toDomainPostDetail(p), nil
}

func (r *Repo) ListPublishedPosts(ctx context.Context, limit, offset int) ([]domain.Post, error) {
	ps, err := r.client.Post.Query().
		Where(post.PublishedEQ(true)).
		Order(ent.Desc(post.FieldCreatedAt)).
		Limit(limit).
		Offset(offset).
		WithAuthor().
		All(ctx)
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
	return r.client.Post.Query().
		Where(post.HasAuthorWith(user.IDEQ(int(authorID)))).
		Count(ctx)
}

func toDomainUser(u *ent.User) *domain.User {
	return &domain.User{ID: int64(u.ID), Name: u.Name, Email: u.Email, CreatedAt: u.CreatedAt}
}

func toDomainPost(p *ent.Post) *domain.Post {
	return &domain.Post{
		ID:        int64(p.ID),
		Title:     p.Title,
		Body:      p.Body,
		AuthorID:  int64(p.Edges.Author.ID),
		Published: p.Published,
		CreatedAt: p.CreatedAt,
	}
}

func toDomainComment(c *ent.Comment) *domain.Comment {
	return &domain.Comment{
		ID:        int64(c.ID),
		PostID:    int64(c.Edges.Post.ID),
		AuthorID:  int64(c.Edges.Author.ID),
		Body:      c.Body,
		CreatedAt: c.CreatedAt,
	}
}

func toDomainPostDetail(p *ent.Post) *domain.PostDetail {
	detail := &domain.PostDetail{
		Post: *toDomainPost(p),
	}

	if a := p.Edges.Author; a != nil {
		detail.Author = *toDomainUser(a)
	}

	for _, t := range p.Edges.Tags {
		detail.Tags = append(detail.Tags, domain.Tag{ID: int64(t.ID), Name: t.Name})
	}

	for _, c := range p.Edges.Comments {
		if c.Edges.Author != nil {
			detail.Comments = append(detail.Comments, *toDomainComment(c))
		}
	}

	return detail
}

func rollback(tx *ent.Tx, err error) error {
	if rerr := tx.Rollback(); rerr != nil {
		err = fmt.Errorf("%w: %v", err, rerr)
	}
	return err
}
