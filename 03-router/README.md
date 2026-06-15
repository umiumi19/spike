# スパイク3: ルーター差し替え(net/http ⇔ chi ⇔ Echo)

目的: **差し替えが本当に数行か** を自分の目で確認 + net/http の基礎学習。

ポイントは、`ops.go` の **オペレーション登録コードは3つとも完全に同じ**で、
変わるのは各 `*_main.go` の「ルーター生成 + アダプタ + 起動」の **3行だけ** だということ。

## 依存
```sh
go get github.com/danielgtaylor/huma/v2
go get github.com/go-chi/chi/v5
go get github.com/labstack/echo/v4
```

## 動かし方
3ファイルは同じ package で main 関数が重複するので、**1つずつ** build タグ等で切り替えるか、
別ディレクトリにコピーして試すのが手軽。まずは差分(下記)を眺めるだけでも目的は果たせる。

## 差分はこれだけ(3行)
- net/http: `mux := http.NewServeMux()` / `humago.New(mux, cfg)` / `http.ListenAndServe(":8888", mux)`
- chi:      `r := chi.NewMux()`         / `humachi.New(r, cfg)`  / `http.ListenAndServe(":8888", r)`
- echo:     `e := echo.New()`           / `humaecho.New(e, cfg)` / `e.Start(":8888")`

→ `huma.Register(...)` の中身(API本体)は一切変わらない = 差し替えコストは軽い、を実感する。
