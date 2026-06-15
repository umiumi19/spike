# go-jet

## セットアップ
```sh
# 1) 共通スキーマを適用済みのDBが必要(jetは既存DBから生成)
# 2) ジェネレータ導入
go install github.com/go-jet/jet/v2/cmd/jet@latest
# 3) DBスキーマから型付きSQLビルダを生成
jet -dsn="postgres://spike:spike@localhost:5432/spike?sslmode=disable" \
    -schema=public -path=./.gen
```
生成後、`.gen/` に table/model パッケージができる。`query_TODO.go` を埋める。
