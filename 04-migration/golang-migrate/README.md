# golang-migrate
```sh
# 導入(mac例)
brew install golang-migrate

DB="postgres://spike:spike@localhost:5432/spike?sslmode=disable"
migrate -path . -database "$DB" up        # 適用
migrate -path . -database "$DB" down 1    # 1つ戻す
# 新規作成:
migrate create -ext sql -dir . add_email_to_users
#   → 000002_add_email_to_users.up.sql / .down.sql が出る → 中身を書く
```
up/down を手で書く素朴さ(差分自動生成なし)を体感する。
