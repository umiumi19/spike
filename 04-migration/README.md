# スパイク4: マイグレーション(Atlas / goose / golang-migrate)

目的: 同じスキーマ変更を3ツールで運用してみて手触りを比べる。
(最終的には ent 統合に決定済みだが、独立ツールの感覚を知っておくと ent の理解も深まる)

共通の「お題」: 既存の users/posts/comments に、users へ `email text` 列を1つ足す。

## 各ツール
- `goose/`          … 連番SQL + Go移行可。シンプル。
- `golang-migrate/` … up/down ペアのSQL。最も素朴。
- `atlas/`          … 宣言的 + lint。高機能だが学習コスト高。

DB接続: `postgres://spike:spike@localhost:5432/spike?sslmode=disable`
