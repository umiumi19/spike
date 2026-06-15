# 技術選定スパイク・リポジトリ(骨組み)

迷った4要素を「手で動かして確かめる」ための最小スパイク集。
ロジックの肝は `TODO` にしてあります。自分で埋めて書き心地・つながりを体感するのが目的です。

## 目的の枠(重要)
- これは「再検討」ではなく **確認・学習**。選んだ道が動くか/手に馴染むかを見る。
- **時間を区切る**(目安: 週末1〜2日)。全部を作り込んで比較沼にはまらない。

## 前提ツール
- Go 1.22+
- Docker(ローカルPostgres用)
- Flutter SDK(2のDart生成のみ)
- 各CLIは各スパイクのREADMEでインストール

## まず共通のDBを起動
```sh
docker compose up -d            # localhost:5432 に postgres 起動
# 接続: postgres://spike:spike@localhost:5432/spike?sslmode=disable
```
共通スキーマ `schema/schema.sql`(users / posts / comments)を使う。
UUID主キー・論理削除(deleted_at)・created_at/updated_at を最初から入れてある(オフライン要件に合わせた本番想定の形)。

## 4つのスパイク
| # | フォルダ | 何を確かめる | 性質 |
|---|----------|--------------|------|
| 1 | `1-db-access/` | ent / go-jet / sqlc で**同じ複雑クエリ**を書き比べ | 書き心地が決め手 |
| 2 | `2-api-toolchain/` | net/http+Huma → OpenAPI → Dart生成が**最後までつながるか** | つながり検証 |
| 3 | `3-router-swap/` | net/http ⇔ chi ⇔ Echo の**差し替えが数行か** | 差し替えコスト確認 |
| 4 | `4-migration/` | Atlas / goose / golang-migrate の**運用感** | 手触り比較 |

## おすすめの順番
1. `docker compose up -d` でDB起動
2. **4(マイグレーション)** で schema を適用 → これで 1 の jet/sqlc が生成元のスキーマを得られる
3. **1(DBアクセス)** で同じクエリを3通り書き比べ ← 本命
4. **2(APIツールチェーン)** で一気通貫を1本通す
5. **3(ルーター差し替え)** を最後に軽く確認

## go.mod について
ルートに最小の `go.mod` を置いてある。各スパイクで使うライブラリは
`go get` で都度追加し、`go mod tidy` で整える(各READMEに記載)。
