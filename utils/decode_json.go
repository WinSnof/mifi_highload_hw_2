package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

func DecodeAndValidate(w http.ResponseWriter, r *http.Request, v Validatable) bool {
	err := json.NewDecoder(r.Body).Decode(v)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxError):
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid JSON format at position %d", syntaxError.Offset))
		case errors.Is(err, io.EOF):
			RespondWithError(w, http.StatusBadRequest, "Request body must not be empty")
		case errors.As(err, &unmarshalTypeError):
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid value for field %q at position %d", unmarshalTypeError.Field, unmarshalTypeError.Offset))
		default:
			RespondWithError(w, http.StatusBadRequest, "Bad Request: "+err.Error())
		}

		return true
	}

	if err := v.Validate(); err != nil {
		RespondWithError(w, http.StatusBadRequest, "Validation error: "+err.Error())
		return true
	}

	return false
}
