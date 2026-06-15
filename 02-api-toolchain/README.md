# スパイク2: APIツールチェーン(net/http + Huma → OpenAPI → Dart生成)

目的は比較ではなく **"最後までつながるか" の検証**。
Goで型を書く → OpenAPIが出る → Dartクライアントが生成される、を一度通す。

## 手順
```sh
# 1) 依存
go get github.com/danielgtaylor/huma/v2

# 2) サーバ起動(main.go の TODO を最小限埋めてから)
go run .
#   → http://localhost:8888/docs に APIドキュメント
#   → http://localhost:8888/openapi.yaml に仕様

# 3) OpenAPI仕様を書き出す
./export-openapi.sh        # openapi.yaml を保存

# 4) Dartクライアントを生成
./gen-dart.sh              # ./client にDartクライアントが出る
```

## 落とし穴チェック
- Huma 既定は OpenAPI 3.1。Dartジェネレータが詰まったら `/openapi-3.0.yaml`(3.0.3版)を使う。
  → `export-openapi.sh` 内のURLを切り替えて試す。
