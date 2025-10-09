package handler

import (
	"net/http"
	"strings"

	"github.com/zakisanbaiman/go-handson01/config"
)

// DefaultCORSOptions デフォルトのCORS設定
func DefaultCORSOptions() *config.CORSOptions {
	return &config.CORSOptions{
		AllowedOrigins: []string{"*"}, // 本番環境では具体的なドメインを指定
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Requested-With",
		},
		MaxAge: 86400, // 24時間
	}
}

// CORSMiddleware CORSミドルウェア
func CORSMiddleware(options *config.CORSOptions) func(next http.Handler) http.Handler {
	if options == nil {
		options = DefaultCORSOptions()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// 許可されたOriginかチェック
			if isAllowedOrigin(origin, options.AllowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
			}

			// プリフライトリクエスト（OPTIONS）の処理
			if r.Method == "OPTIONS" {
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(options.AllowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(options.AllowedHeaders, ", "))
				w.Header().Set("Access-Control-Max-Age", string(rune(options.MaxAge)))
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.WriteHeader(http.StatusOK)
				return
			}

			// 通常のリクエスト
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")

			next.ServeHTTP(w, r)
		})
	}
}

// isAllowedOrigin 許可されたOriginかチェック
func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}

	for _, allowed := range allowedOrigins {
		if allowed == "*" {
			return true
		}
		if allowed == origin {
			return true
		}
		// ワイルドカードサブドメインのサポート
		if strings.HasPrefix(allowed, "*.") {
			domain := strings.TrimPrefix(allowed, "*.")
			if strings.HasSuffix(origin, domain) {
				return true
			}
		}
	}
	return false
}
