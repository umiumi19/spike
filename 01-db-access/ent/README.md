# ent

## セットアップ
```sh
go get entgo.io/ent
go get github.com/lib/pq   # PostgreSQL ドライバ(ent 本体には含まれない)

# エンティティの雛形を生成
go run -mod=mod entgo.io/ent/cmd/ent new User
go run -mod=mod entgo.io/ent/cmd/ent new Car Group
# → ent/schema/{user.go, car.go, group.go} ができる
```

> **注意**: `go run -mod=mod entgo.io/ent/cmd/ent` 実行時に
> `github.com/clipperhouse/displaywidth` のコンパイルエラーが出る場合は
> バージョンを固定する。
> ```sh
> go get github.com/clipperhouse/displaywidth@v0.11.0
> ```
> `go mod tidy` で消えるので消えたら再実行する。

## スキーマ定義

`ent/schema/*.go` にフィールドと edge(関係) を定義する。

**User**
- fields: `age` (int, positive), `name` (string, default "unknown")
- edges: `cars` → Car (1対多), `groups` ← Group (多対多・逆参照)

**Car**
- fields: `model` (string), `registered_at` (time)
- edges: `owner` ← User (逆参照, Unique → 1台のオーナーは1人)

**Group**
- fields: `name` (string, regexp `[a-zA-Z_]+`)
- edges: `users` → User (多対多 → `group_users` 中間テーブルが生成される)

## コード生成

```sh
go generate ./ent
# ent/generate.go の //go:generate が呼ばれ、型付きクライアントが生成される
```

## マイグレーション

### 方法1: Schema.Create (開発用・自動)

```go
client.Schema.Create(context.Background())
```

コード内で呼ぶだけで DB にテーブルが作られる。差分管理はされない。

### 方法2: atlas migrate (推奨)

```sh
# 差分 SQL を生成
atlas migrate diff migration_name \
  --dir "file://ent/migrate/migrations" \
  --to "ent://ent/schema" \
  --dev-url "docker://postgres/16/test?search_path=public"

# DB に適用
atlas migrate apply \
  --dir "file://ent/migrate/migrations" \
  --url "postgres://spike:spike@localhost:5432/spike?sslmode=disable"
```

`01-db-access/` ディレクトリから実行する。

## DB リセット

```sh
# 確実な方法(ボリュームごと削除)
docker compose down -v && docker compose up -d

# atlas で削除する場合
atlas schema clean --url "postgres://spike:spike@localhost:5432/spike?sslmode=disable"
# ⚠️ public スキーマ自体も消えるので、atlas migrate apply 前に再作成が必要
psql "postgres://spike:spike@localhost:5432/spike?sslmode=disable" -c "CREATE SCHEMA public;"
```

## 生成されるテーブル

| テーブル | 概要 |
|---|---|
| `users` | User エンティティ |
| `cars` | Car エンティティ (`user_cars` FK で users を参照) |
| `groups` | Group エンティティ |
| `group_users` | Group-User 多対多の中間テーブル |
