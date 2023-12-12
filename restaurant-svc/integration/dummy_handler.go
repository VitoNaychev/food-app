package integration

import (
	"errors"
	"net/http"

	"github.com/VitoNaychev/food-app/httperrors"
)

func DummyHandlerFunc(w http.ResponseWriter, r *http.Request) {
	httperrors.WriteJSONError(w, http.StatusMisdirectedRequest, errors.New("called dummy handler, did you really want to?"))
}

var DummyHandler = http.HandlerFunc(DummyHandlerFunc)
