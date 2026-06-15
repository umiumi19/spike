# sqlc

## セットアップ
```sh
# 1) 共通スキーマを適用(スパイク4のいずれかで or 直接psqlで)
# 2) sqlc 導入
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
# 3) コード生成
sqlc generate     # query.sql + schema.sql から ./db に型付きコードを生成
```
`query.sql` の TODO を埋めてから generate すると、型付き関数が ./db に出る。
動的条件は sqlc 単体だと苦しい点を、ここで実感できるはず。
