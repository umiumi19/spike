# db-compare — gorm / ent / sqlc を同一ドメインで比較

同じ「ブログ」ドメイン（User / Post / Tag / Comment）に対して、同じ
`domain.Repository` インターフェースを **gorm・ent・sqlc の3通りで実装**し、
同一シナリオ（`scenario/scenario.go`）を流して挙動・書き味・生成物を比較します。

DB は Docker の Postgres を1コンテナだけ立て、ツールごとに別データベース
（`blog_gorm` / `blog_ent` / `blog_sqlc`）を使います。

```
domain/           共通モデルと Repository インターフェース（依存ゼロ）
scenario/         3実装に共通で流すシナリオ
gorm/             gorm 実装 + main
  model/          GORM モデル定義（構造体タグでスキーマ定義）
  gen/            gorm-gen コード生成スクリプト
  query/          gorm-gen 生成コード
ent/              ent 実装 + main
  ent/            ent スキーマ(schema/) + 生成コード（go generate で生成）
  migrate/        Atlas マイグレーション生成スクリプト（main.go）
  migrations/     生成済みマイグレーションファイル（golang-migrate 形式）
sqlc/             sqlc 実装 + main
  gen/            sqlc 生成コード（sqlc generate で生成）
  migrations/     生成済みマイグレーションファイル（Atlas ネイティブ形式）
  schema.sql      スキーマ定義
  query.sql       クエリ定義
  sqlc.yaml       sqlc 設定
  atlas.hcl       Atlas 設定
```

---

## 既存コードを動かす

### 0. 前提

- Go 1.22+
- Docker / Docker Compose
- 追加ツール: `sqlc`, `atlas`

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
brew install ariga/tap/atlas
```

### 1. Postgres を起動

```bash
docker compose up -d
docker exec -it dbcompare-pg psql -U blog -c '\l' | grep blog_
```

`blog_gorm` `blog_ent` `blog_sqlc` が見えればOK。
（DB は初回起動時にだけ作られます。作り直したいときは `docker compose down -v`）

### 2. 依存を入れる

```bash
go mod tidy
```

### 3. gorm を動かす

```bash
go run ./gorm
```

### 4. ent を動かす

```bash
go generate ./ent/ent   # 生成コードを作る
go run ./ent
```

### 5. sqlc を動かす

```bash
cd sqlc
sqlc generate                        # 生成コードを作る
atlas migrate apply --env local      # schema.sql をDBに適用
cd ..
go run ./sqlc
```

### 6. 3つの出力が一致することを確認

```bash
go run ./gorm
go run ./ent
go run ./sqlc
```

いずれも user / post / detail(tags=2, comments=1) / count などが
同じ内容で出れば、3実装が等価に動いています。

---

## 1から作り直す手順

各ツールで「何を書いて、いつコマンドを実行するか」の流れをまとめます。

### gorm

**書くもの:**

1. `gorm/model/model.go` — Go 構造体に gorm タグを付けてスキーマを表現する
   - フィールドタグ（`gorm:"primaryKey"`, `gorm:"uniqueIndex"` など）
   - リレーションタグ（`gorm:"foreignKey:AuthorID"`, `gorm:"many2many:post_tags;"` など）
2. `gorm/repository.go` — `domain.Repository` を実装する
3. `gorm/main.go` — `repo.Migrate()`（内部で `AutoMigrate`）を呼んでから `scenario.Run`

**コマンド実行のタイミング:**

```
構造体を書く
  → go run ./gorm   # AutoMigrate がテーブルを作り、そのまま実行まで通る
```

> スキーマを変更したい場合: 構造体を修正して再実行するだけ。`AutoMigrate` が差分を自動適用する。

---

### ent

**書くもの:**

1. `ent/ent/schema/*.go` — エンティティごとにスキーマファイルを書く
   - `Fields()` でカラムを定義（`field.String("name")`, `field.Bool("published")` など）
   - `Edges()` でリレーションを定義（`edge.To("posts", Post.Type)` など）
2. `ent/ent/entc.go` — Atlas 連携の設定（`WithAtlasHCL` などのオプション）
3. `ent/repository.go` — 生成された型安全な Builder API を使って `domain.Repository` を実装する
4. `ent/main.go` — `client.Schema.Create(ctx)` を呼んでから `scenario.Run`
5. `ent/migrate/main.go` — Atlas でマイグレーションファイルを生成するスクリプト（バージョン管理したい場合）

**コマンド実行のタイミング:**

```
schema/*.go を書く
  → go generate ./ent/ent    # ent/ent/ 配下に大量のコードが生成される
  → repository.go を書く     # 生成された Client/Builder を使って実装する
  → go run ./ent             # Schema.Create がテーブルを作り、そのまま実行まで通る
```

> スキーマを変更したい場合:
> ```
> schema/*.go を修正
>   → go generate ./ent/ent
>   → go run ./ent   # Schema.Create が差分を自動適用
> ```

> Atlas でマイグレーションファイルを残したい場合（任意）:
> ```
> go generate ./ent/ent
>   → go run ent/migrate/main.go <マイグレーション名>
>      # ent/migrations/ に .up.sql / .down.sql が生成される
> ```

---

### sqlc

**書くもの:**

1. `sqlc/schema.sql` — CREATE TABLE 文でスキーマを定義する
2. `sqlc/query.sql` — `-- name: QueryName :one/:many/:exec` アノテーション付きで SQL クエリを書く
3. `sqlc/sqlc.yaml` — エンジン・入力ファイル・出力先・パッケージ名を指定する設定ファイル
4. `sqlc/atlas.hcl` — Atlas の環境設定（接続先 URL・dev URL・migrations ディレクトリ）
5. `sqlc/repository.go` — `gen/` の生成コードを使って `domain.Repository` を実装する
6. `sqlc/main.go` — DB 接続を作って `scenario.Run` に渡す

**コマンド実行のタイミング:**

```
schema.sql と query.sql と sqlc.yaml を書く
  → cd sqlc && sqlc generate    # gen/ に models.go / query.sql.go / db.go が生成される
  → repository.go を書く        # 生成された Queries 型とパラメータ型を使って実装する
  → main.go を書く

atlas.hcl を書く
  → atlas migrate diff initial --env local
     # migrations/ に初回マイグレーション SQL が生成される
  → atlas migrate apply --env local
     # blog_sqlc にテーブルが作られる
  → go run ./sqlc
```

> スキーマを変更したい場合:
> ```
> schema.sql を修正
>   → sqlc generate                           # gen/ の型を更新
>   → atlas migrate diff <変更名> --env local  # 差分 SQL を生成
>   → atlas migrate apply --env local         # DB に適用
>   → go run ./sqlc
> ```

> `atlas migrate diff` の仕組み:
> `docker://postgres/16/dev` の一時コンテナに既存マイグレーションを replay して
> 「現在の DB 状態」を導出し、`schema.sql` との差分だけを新ファイルとして出力する。

---

## トラブルシュート

- **`go run ./ent` がビルドできない**: `go generate ./ent/ent` 未実行。`ent/ent/` が空。
- **`go run ./sqlc` がビルドできない**: `sqlc generate` 未実行で `sqlc/gen/` が空。
- **sqlc で `relation "users" does not exist`**: `atlas migrate apply --env local` 未実行。
- **接続できない**: `docker compose ps` で db が healthy か確認。ポート 5432 衝突に注意。
- **ent の並び替えでコンパイルエラー**: ent のバージョンで Order API が変わります。
  本コードは `ent.Desc(post.FieldCreatedAt)` を使用。新しめの版なら
  `post.ByCreatedAt(sql.OrderDesc())`（`entgo.io/ent/dialect/sql` を import）に置き換えてください。
