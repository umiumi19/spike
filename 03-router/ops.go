package routerswap

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

// ★ ここはルーターに依存しない。3つのどのアダプタでも同じこのコードを使う。
type PingOutput struct {
	Body struct {
		Message string `json:"message"`
	}
}

func RegisterOps(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "ping",
		Method:      http.MethodGet,
		Path:        "/ping",
	}, func(ctx context.Context, _ *struct{}) (*PingOutput, error) {
		out := &PingOutput{}
		out.Body.Message = "pong"
		return out, nil
	})
}
