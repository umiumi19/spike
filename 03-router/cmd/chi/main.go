package main

import (
	"encoding/json"
	"log"
	"net/http"
	"spike/03-router/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


// ============================================================
// ミドルウェア
// chi は http.Handler ラッパーをそのまま使う — net/http 資産が流用できる。
// ============================================================

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Auth") != "secret" {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}
		next.ServeHTTP(w, r)
	})
}

// ============================================================
// ルート登録
// r.Route でグループ化 → /api/v1 プレフィックスを一カ所で管理。
// r.Group でサブグループを作り、そこだけ authMiddleware を追加。
// ============================================================

func main() {
	r := chi.NewRouter()

	// 組み込みミドルウェアを宣言的にスタック
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/items", listItems)
		r.Get("items/{id}", getItem)

		// このサブグループにだけ auth を適用
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)
			r.Post("/items", createItem)
		})
	})

	log.Println("chi  :8002")
	log.Fatal(http.ListenAndServe(":8002", r))
}

// ============================================================
// ハンドラー: func(http.ResponseWriter, *http.Request) — net/http と同じ型
// パスパラメータ: chi.URLParam(r, "id")
// ============================================================


func health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func listItems(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, store.List())
}

func getItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id") // chi 方式
	item, ok := store.Get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	item := store.Create(body.Name)
	writeJSON(w, http.StatusCreated, item)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
