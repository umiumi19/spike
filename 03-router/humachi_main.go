//go:build humachi

package routerswap

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

// 実行: go run -tags humachi .
func main() {
	r := chi.NewMux()                                          // ①
	api := humachi.New(r, huma.DefaultConfig("swap", "1.0.0")) // ②
	RegisterOps(api)                                           // 共通(不変)
	http.ListenAndServe(":8888", r)                            // ③
}
