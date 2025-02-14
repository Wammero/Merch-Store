package jwt

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

// Ключ для хранения данных в контексте
type contextKey string

const UserContextKey = contextKey("user")

func SetSecret(secret string) {
	jwtSecret = []byte(secret)
}

func JWTValidator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Неавторизован.", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Неавторизован.", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		claims := &jwt.RegisteredClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Неавторизован.", http.StatusUnauthorized)
			return
		}

		// Сохраняем имя пользователя в контекст
		ctx := context.WithValue(r.Context(), UserContextKey, claims.Subject)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		Issuer:    "merch-store",
		Subject:   username,
	})

	return token.SignedString(jwtSecret)
}
