package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// errorResponse — стандартный JSON-ответ с ошибкой.
type errorResponse struct {
	Error string `json:"error"`
}

// writeJSON сериализует data в JSON и пишет в w.
func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	b, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err.Error())))
		return
	}
	w.Write(b)
}

// itoa конвертирует int64 в строку без импорта strconv в каждом файле.
func itoa(n int64) string {
	return fmt.Sprintf("%d", n)
}
