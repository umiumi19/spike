# Atlas
```sh
# 導入
curl -sSf https://atlasgo.sh | sh

# 宣言的: 望むスキーマ(schema.sql)と現状の差分からマイグレーションを自動生成
atlas migrate diff --env local add_email_to_users
#   → migrations/ に差分SQLが生成される(ここが goose/golang-migrate との違い)

# lint(危険な変更の検知)
atlas migrate lint --env local --latest 1

# 適用
atlas migrate apply --env local
```
TODO: schema.sql の users に `email text` を足してから `migrate diff` を実行し、
差分が自動生成される体験と lint の警告を見る。
