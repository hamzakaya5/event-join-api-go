package middleware

import (
	"context"
	"gamegos_case/database"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type CtxUserKey struct{}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "missing Authorization header", http.StatusUnauthorized)
			return
		}
		parts := strings.Fields(auth)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "invalid Authorization header", http.StatusUnauthorized)
			return
		}
		token := parts[1]
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		val, err := database.RedisClient.Get(ctx, token).Result()
		if err == redis.Nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Printf("redis error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxUserKey{}, val)))
	}
}
