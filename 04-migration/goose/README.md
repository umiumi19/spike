# goose
```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="postgres://spike:spike@localhost:5432/spike?sslmode=disable"

goose -dir . up        # 適用
goose -dir . status    # 状態確認
goose -dir . down      # 1つ戻す
# 新規作成:
goose -dir . create add_email_to_users sql   # 00002_add_email_to_users.sql が出る → 中身を書く
```
TODO: 00001 にスキーマ作成、00002 で `ALTER TABLE users ADD COLUMN email text;` を書いて up/down を体感。
