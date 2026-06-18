# スパイク2: APIツールチェーン (Huma + oapi-codegen)

Huma 公式チュートリアル ([huma.rocks](https://huma.rocks/tutorial/installation/)) に沿って、Go の struct から OpenAPI を自動生成し、そこから型安全な Go クライアントを生成するフローを通す。

## 構成

```
.
├── main.go          # サーバー本体 (Huma v2 + net/http)
├── main_test.go     # humatest によるハンドラーテスト
├── openapi.yaml     # 生成済み OpenAPI 3.0.3 仕様
├── sdk/sdk.go       # oapi-codegen で生成した Go クライアント
└── client/client.go # SDK を使った呼び出し例
```

## エンドポイント

| Method | Path | 概要 |
|--------|------|------|
| GET | `/greeting/{name}` | 挨拶を返す |
| POST | `/reviews` | レビューを登録する (未永続化) |

## 使い方

### サーバー起動

```sh
go run . [--port 8888]
```

- `http://localhost:8888/docs` — Swagger UI
- `http://localhost:8888/openapi.json` — OpenAPI 仕様 (JSON)

### OpenAPI 仕様の書き出し

```sh
# OpenAPI 3.0.3 (YAML) を標準出力に出力
go run . openapi
```

### テスト実行

```sh
go test ./...
```

`humatest` を使い、HTTP サーバーを立てずにハンドラーを直接テストする。

### Go クライアントの生成 (oapi-codegen)

```sh
# oapi-codegen がなければインストール
go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

# sdk/sdk.go を再生成
oapi-codegen -package sdk -generate types,client openapi.yaml > sdk/sdk.go
```

### クライアント例の実行

サーバーを起動した状態で:

```sh
go run ./client
# → Hello, world!
```

## メモ

- Huma のデフォルトは OpenAPI 3.1。oapi-codegen が 3.1 を完全サポートしていないため、CLI の `DowngradeYAML()` で 3.0.3 に変換してから生成する。
- 入力バリデーション (`maxLength`, `minimum`, `maximum`) は Huma が JSON Schema 経由でハンドラー呼び出し前に自動検証する。
