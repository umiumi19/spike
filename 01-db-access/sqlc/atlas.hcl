env "local" {
  src = "file://schema.sql"
  url = "postgres://blog:blog@localhost:5432/blog_sqlc?sslmode=disable"
  dev = "docker://postgres/16/dev"
  migration {
    dir = "file://migrations"
  }
}
