package main

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const ctxUserID = contextKey("user_id")

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := r.Header.Get("Authorization")
		if h == "" {
			http.Error(w, "missing auth", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid auth", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]
		tok, err := jwt.Parse(tokenStr, func(tok *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		}, jwt.WithValidMethods([]string{"HS256"}))
		if err != nil || !tok.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		claims := tok.Claims.(jwt.MapClaims)
		sub := int(claims["sub"].(float64))
		ctx := context.WithValue(r.Context(), ctxUserID, sub)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}