# スパイク3: ルーター比較 (net/http / chi / echo)

同じ API を3つのルーターで実装し、書き味の違いを実感する。

## 構成

```
03-router/
├── store/store.go               # 共通のインメモリストア
├── cmd/nethttprouter/main.go    # 標準ライブラリ net/http
├── cmd/chirouter/main.go        # chi v5
└── cmd/echorouter/main.go       # echo v4
```

## 実装するAPI (3つとも同じ)

| Method | Path | 認証 | 概要 |
|--------|------|------|------|
| GET | `/health` | - | ヘルスチェック |
| GET | `/api/v1/items` | - | 一覧取得 |
| GET | `/api/v1/items/{id}` | - | ID 指定取得 |
| POST | `/api/v1/items` | `X-Auth: secret` | 新規作成 |

## 起動

```sh
# spike/ ルートから実行
go run ./03-router/cmd/nethttprouter   # :8001
go run ./03-router/cmd/chirouter      # :8002
go run ./03-router/cmd/echorouter     # :8003
```

## 動作確認

```sh
# 一覧
curl http://localhost:800X/api/v1/items

# ID 指定
curl http://localhost:800X/api/v1/items/1

# 作成 (認証あり)
curl -X POST http://localhost:800X/api/v1/items \
  -H "X-Auth: secret" \
  -H "Content-Type: application/json" \
  -d '{"name":"cherry"}'

# 認証なし → 401
curl -X POST http://localhost:800X/api/v1/items \
  -H "Content-Type: application/json" \
  -d '{"name":"cherry"}'

# 存在しないルートへ GET → chi/echo は 405 自動返却、net/http は 404
curl -X DELETE http://localhost:800X/api/v1/items/1
```

---

## 差分ポイント

### 1. ルートグループ化

```go
// net/http — グループなし。プレフィックスを全ルートに手書き
mux.HandleFunc("GET /api/v1/items", listItems)
mux.HandleFunc("GET /api/v1/items/{id}", getItem)
mux.Handle("POST /api/v1/items", auth(http.HandlerFunc(createItem)))

// chi — r.Route でネスト
r.Route("/api/v1", func(r chi.Router) {
    r.Get("/items", listItems)
    r.Get("/items/{id}", getItem)
    r.Group(func(r chi.Router) {
        r.Use(authMiddleware)
        r.Post("/items", createItem)
    })
})

// echo — e.Group でプレフィックス
v1 := e.Group("/api/v1")
v1.GET("/items", listItems)
v1.GET("/items/:id", getItem)
v1.POST("/items", createItem, authMiddleware)
```

### 2. ミドルウェアの適用

```go
// net/http — ルーター全体 or 1ルートずつ手動ラップ
http.ListenAndServe(":8001", logging(mux))          // 全体
mux.Handle("POST /...", auth(http.HandlerFunc(fn))) // 1ルートだけ

// chi — r.Use() で宣言的にスタック
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)

// echo — e.Use() で同様
e.Use(middleware.Logger())
e.Use(middleware.Recover())
```

### 3. パスパラメータの取得

```go
r.PathValue("id")    // net/http (Go 1.22+)
chi.URLParam(r, "id") // chi
c.Param("id")        // echo
```

### 4. ハンドラーの型

```go
func(w http.ResponseWriter, r *http.Request)  // net/http & chi — 共通
func(c echo.Context) error                    // echo — 非互換
```

echo はハンドラーが `error` を返すため、chi/net/http の既存ハンドラーをそのまま流用できない。
その代わり `c.JSON(code, v)` で Content-Type セットと JSON エンコードが1行になる。

### 5. 未定義メソッドへのレスポンス

登録していない `DELETE /api/v1/items/1` を叩くと:

| | レスポンス |
|--|--|
| net/http | 404 Not Found |
| chi | **405 Method Not Allowed** |
| echo | **405 Method Not Allowed** |

---

## Huma を使う場合はどう変わるか

Huma (02-api-toolchain) はルーターの上にアダプター層を置くフレームワーク。  
**ハンドラーレベルの差異はほぼ吸収される**が、**ミドルウェアとルートグループは依然として router 次第**。

### Huma が吸収してくれること

| このスパイクで確認した違い | Huma での扱い |
|---|---|
| パスパラメータの取得方法 | struct タグで宣言 → Huma が自動で詰める |
| JSON エンコード/デコード | 自動 |
| バリデーション (`maxLength` 等) | struct タグで自動 |
| ハンドラーの型 | 常に `func(ctx, *Input) (*Output, error)` |

`huma.Register()` を使う限り、中身のハンドラーコードはどのルーターを選んでも**一切変わらない**。

### 依然として router が担うこと

ミドルウェアとルートグループは Huma の外側の話なので、router の違いがそのまま出る。

```go
// 認証ミドルウェアを /admin 以下にだけかけたい場合

// chi なら r.Route + r.Use で自然に書ける
r.Route("/admin", func(r chi.Router) {
    r.Use(authMiddleware)
    huma.Register(api, ...)
})

// net/http にはグループ機能がないので工夫が要る
```

ログ出力・CORS・リカバリー・認証といった横断的な処理は Huma が関知しないため、  
chi/echo の `r.Use()` の恩恵がそのまま効く。
