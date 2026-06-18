package main

import (
	"net/http"
	"spike/03-router/store"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// ============================================================
// ミドルウェア
// echo.HandlerFunc を受け取り echo.HandlerFunc を返す — net/http と非互換。
// ルート単位で第3引数以降に渡すこともできる（createItem 参照）。
// ============================================================

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get("X-Auth") != "secret" {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		}
		return next(c)
	}
}

func main() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", health)

	v1 := e.Group("/api/v1")
	v1.GET("/items", listItems)
	v1.GET("/items/:id", getItem)              // ← ":id" 記法
	v1.POST("/items", createItem, authMiddleware) // ← ルート単位でミドルウェアを追加

	e.Logger.Fatal(e.Start(":8003"))
}

// ============================================================
// ハンドラー: func(echo.Context) error — net/http と非互換
// パスパラメータ: c.Param("id")
// JSONレスポンス: c.JSON(code, v) で Content-Type も自動セット
// エラーは return するだけ — echo がエラーレスポンスに変換
// ============================================================

func health(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}

func listItems(c echo.Context) error {
	return c.JSON(http.StatusOK, store.List())
}

func getItem(c echo.Context) error {
	id := c.Param("id") // echo 方式
	item, ok := store.Get(id)
	if !ok {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}
	return c.JSON(http.StatusOK, item)
}

func createItem(c echo.Context) error {
	var body struct {
		Name string `json:"name"`
	}
	if err := c.Bind(&body); err != nil { // echo は Bind でデコード
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	item := store.Create(body.Name)
	return c.JSON(http.StatusCreated, item)
}
