-- +goose Up
SELECT 'up SQL query';
ALTER TABLE "users" ADD COLUMN "email" text NULL;

-- +goose Down
SELECT 'down SQL query';
 ALTER TABLE users DROP COLUMN IF EXISTS email;