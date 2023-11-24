package errorresponse

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
)

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Details string `json:"details"`
}

func WriteJSONError(w http.ResponseWriter, statusCode int, err error) {
	errorResponse := ErrorResponse{
		Error:   reflect.TypeOf(err).Name(),
		Message: err.Error(),
	}

	wrappedError := errors.Unwrap(err)
	if wrappedError != nil {
		errorResponse.Details = wrappedError.Error()
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(errorResponse)
}
