-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  "email" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "users_email_key" UNIQUE ("email")
);
-- Create "posts" table
CREATE TABLE "public"."posts" (
  "id" bigserial NOT NULL,
  "title" text NOT NULL,
  "body" text NOT NULL,
  "author_id" bigint NOT NULL,
  "published" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "posts_author_id_fkey" FOREIGN KEY ("author_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create "comments" table
CREATE TABLE "public"."comments" (
  "id" bigserial NOT NULL,
  "post_id" bigint NOT NULL,
  "author_id" bigint NOT NULL,
  "body" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "comments_author_id_fkey" FOREIGN KEY ("author_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "comments_post_id_fkey" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "tags" table
CREATE TABLE "public"."tags" (
  "id" bigserial NOT NULL,
  "name" text NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "tags_name_key" UNIQUE ("name")
);
-- Create "post_tags" table
CREATE TABLE "public"."post_tags" (
  "post_id" bigint NOT NULL,
  "tag_id" bigint NOT NULL,
  PRIMARY KEY ("post_id", "tag_id"),
  CONSTRAINT "post_tags_post_id_fkey" FOREIGN KEY ("post_id") REFERENCES "public"."posts" ("id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "post_tags_tag_id_fkey" FOREIGN KEY ("tag_id") REFERENCES "public"."tags" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
