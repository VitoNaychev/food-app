package httperrors

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSONError(w http.ResponseWriter, statusCode int, err error) {
	errorResponse := ErrorResponse{
		Error: err.Error(),
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}

func HandleBadRequest(w http.ResponseWriter, err error) {
	WriteJSONError(w, http.StatusBadRequest, err)
}

func HandleInternalServerError(w http.ResponseWriter, err error) {
	WriteJSONError(w, http.StatusInternalServerError, err)
}

func HandleNotFound(w http.ResponseWriter, err error) {
	WriteJSONError(w, http.StatusNotFound, err)
}

func HandleUnauthorized(w http.ResponseWriter, err error) {
	WriteJSONError(w, http.StatusUnauthorized, err)
}
