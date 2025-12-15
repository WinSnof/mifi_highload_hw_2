package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Validatable interface {
	Validate() error
}
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errResp := ErrorResponse{
		Status:  code,
		Message: message,
	}

	// Логируем ошибку, если не удалось закодировать ответ об ошибке
	if err := json.NewEncoder(w).Encode(errResp); err != nil {
		fmt.Printf("CRITICAL ERROR: Failed to encode error response: %v\n", err)
	}
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		fmt.Printf("CRITICAL ERROR: Failed to encode successful response: %v\n", err)
	}
}
