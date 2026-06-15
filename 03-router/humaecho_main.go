//go:build humaecho

package routerswap

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humaecho"
	"github.com/labstack/echo/v4"
)

// 実行: go run -tags humaecho .
func main() {
	e := echo.New()                                             // ①
	api := humaecho.New(e, huma.DefaultConfig("swap", "1.0.0")) // ②
	RegisterOps(api)                                            // 共通(不変)
	e.Start(":8888")                                            // ③
}
