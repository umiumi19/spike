package main

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

// レスポンス型(Goの型がそのままOpenAPIスキーマになる)
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ListItemsOutput struct {
	Body []Item `json:"body"`
}

func main() {
	mux := http.NewServeMux()
	api := humago.New(mux, huma.DefaultConfig("Spike API", "1.0.0"))

	// 1オペレーションだけ登録(型からOpenAPIが自動生成される)
	huma.Register(api, huma.Operation{
		OperationID: "list-items",
		Method:      http.MethodGet,
		Path:        "/items",
		Summary:     "アイテム一覧",
	}, func(ctx context.Context, _ *struct{}) (*ListItemsOutput, error) {
		// TODO: 本番ではここでDB(ent)から取得する。今は動作確認用の固定データ。
		return &ListItemsOutput{Body: []Item{
			{ID: "1", Name: "sample"},
		}}, nil
	})

	// TODO: 入力(パスパラメータ/ボディ)付きのオペレーションも1つ足して、
	//       生成されるDartクライアントの型を確認すると学びが大きい。

	http.ListenAndServe(":8888", mux)
}
