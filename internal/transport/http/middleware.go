package http

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

func (h *Handler) JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			h.error(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}
		tokenString = tokenString[7:]

		jwtSecret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			h.log.Error("failed to parse token: ", "error", err.Error())
			h.error(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			h.log.Error("failed to parse token claims")
			h.error(w, http.StatusUnauthorized, fmt.Errorf("invalid token"))
			return
		}
		ctx := context.WithValue(r.Context(), "userID", claim["sub"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
