-- Runs automatically the first time the Postgres volume is initialized.
-- Each tool gets its own database so their schema management never collides.
CREATE DATABASE blog_gorm;
CREATE DATABASE blog_ent;
CREATE DATABASE blog_sqlc;
