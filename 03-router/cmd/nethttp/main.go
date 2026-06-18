package main

import (
	"encoding/json"
	"log"
	"net/http"
	"spike/03-router/store"
	"time"
)


// ============================================================
// ミドルウェア
// net/http には Use() がない。
// 適用したいハンドラーを関数でラップするだけ。
// ============================================================

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

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
// グループ化の仕組みはない。
// /api/v1 プレフィックスを全ルートに手動で付ける。
// ============================================================

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", health)
	mux.HandleFunc("GET /api/v1/items", listItems)
	mux.HandleFunc("GET /api/v1/items/{id}", getItem)
	// POST /items だけに auth を適用 — 1ルートずつ個別ラップ
	mux.Handle("POST /api/v1/items", authMiddleware(http.HandlerFunc(createItem)))

	// logging はルーター全体にかける
	log.Println("net/http  :8001")
	log.Fatal(http.ListenAndServe(":8001", logging(mux)))
}

// ============================================================
// ハンドラー: func(http.ResponseWriter, *http.Request)
// パスパラメータ: r.PathValue("id")  — Go 1.22 で追加
// JSONレスポンス: Content-Type を毎回セットする必要がある
// ============================================================

func health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func listItems(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, store.List())
}

func getItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
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
	writeJSON(w, http.StatusOK, item)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

