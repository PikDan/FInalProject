package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5" // go get github.com/golang-jwt/jwt/v5
)

// jwtSecret — секрет для подписи JWT, берём из пароля чтобы смена пароля инвалидировала токены
func jwtSecret() []byte {
	return []byte(os.Getenv("TODO_PASSWORD"))
}

// passwordHash возвращает SHA-256 хэш пароля в виде hex-строки.
// Хэш хранится внутри JWT как claim — по нему проверяем валидность токена.
func passwordHash(pass string) string {
	h := sha256.Sum256([]byte(pass))
	return fmt.Sprintf("%x", h)
}

type signinRequest struct {
	Password string `json:"password"`
}

// signinHandler обрабатывает POST /api/signin.
// Сверяет пароль с TODO_PASSWORD, при совпадении возвращает JWT-токен.
func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req signinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: err.Error()})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")
	if pass == "" {
		// Пароль не задан — аутентификация отключена
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "аутентификация не настроена"})
		return
	}

	if req.Password != pass {
		writeJSON(w, http.StatusUnauthorized, errorResponse{Error: "неверный пароль"})
		return
	}

	// Формируем JWT: в payload кладём хэш пароля
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"hash": passwordHash(pass),
	})

	tokenStr, err := token.SignedString(jwtSecret())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": tokenStr})
}
