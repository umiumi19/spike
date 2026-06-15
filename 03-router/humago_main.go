//go:build humago

package routerswap

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
)

// 実行: go run -tags humago .
func main() {
	mux := http.NewServeMux()                                   // ①
	api := humago.New(mux, huma.DefaultConfig("swap", "1.0.0")) // ②
	RegisterOps(api)                                            // 共通(不変)
	http.ListenAndServe(":8888", mux)                           // ③
}
