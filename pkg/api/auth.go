package api

import (
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// getToken пытается получить JWT из куки "token",
// затем из заголовка Authorization: Bearer <token>
func getToken(r *http.Request) string {
	// Сначала пробуем куку
	if cookie, err := r.Cookie("token"); err == nil && cookie.Value != "" {
		return cookie.Value
	}
	// Затем заголовок Authorization: Bearer <token>
	if header := r.Header.Get("Authorization"); len(header) > 7 && header[:7] == "Bearer " {
		return header[7:]
	}
	return ""
}

func authError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(`{"error":"authentification required"}`))
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")

		// Аутентификация отключена — пропускаем
		if pass == "" {
			next(w, r)
			return
		}

		tokenStr := getToken(r)
		if tokenStr == "" {
			authError(w)
			return
		}

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtSecret(), nil
		})
		if err != nil || !token.Valid {
			authError(w)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			authError(w)
			return
		}

		hash, ok := claims["hash"].(string)
		if !ok || hash != passwordHash(pass) {
			authError(w)
			return
		}

		next(w, r)
	}
}
