-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES ($1, $2)
RETURNING id,
    name,
    email,
    created_at;
-- name: GetUser :one
SELECT id,
    name,
    email,
    created_at
FROM users
WHERE id = $1;
-- name: CreatePost :one
INSERT INTO posts (title, body, author_id, published)
VALUES ($1, $2, $3, $4)
RETURNING id,
    title,
    body,
    author_id,
    published,
    created_at;
-- name: GetPost :one
SELECT id,
    title,
    body,
    author_id,
    published,
    created_at
FROM posts
WHERE id = $1;
-- name: GetOrCreateTag :one
INSERT INTO tags (name)
VALUES ($1) ON CONFLICT (name) DO
UPDATE
SET name = EXCLUDED.name
RETURNING id,
    name;
-- name: AttachTag :exec
INSERT INTO post_tags (post_id, tag_id)
VALUES ($1, $2) ON CONFLICT DO NOTHING;
-- name: AddComment :one
INSERT INTO comments (post_id, author_id, body)
VALUES ($1, $2, $3)
RETURNING id,
    post_id,
    author_id,
    body,
    created_at;
-- name: ListTagsByPost :many
SELECT t.id,
    t.name
FROM tags t
    JOIN post_tags pt ON pt.tag_id = t.id
WHERE pt.post_id = $1
ORDER BY t.name;
-- name: ListCommentByPost :many
SELECT id,
    post_id,
    author_id,
    body,
    created_at
FROM comments
WHERE post_id = $1
ORDER BY created_at;
-- name: ListPublishedPosts :many
SELECT id,
    title,
    body,
    author_id,
    published,
    created_at
FROM posts
WHERE published = TRUE
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
-- name: CountPostsByAuthor :one
SELECT count(*)
FROM posts
WHERE author_id = $1;