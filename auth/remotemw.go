package auth

import (
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/msgtypes"
)

type VerifyJWTFunc func(token string) (msgtypes.AuthResponse, error)

func RemoteAuthenticationMW(handler func(w http.ResponseWriter, r *http.Request), verifyJWT VerifyJWTFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if tokenHeader := r.Header.Get("Token"); tokenHeader == "" {
			httperrors.WriteJSONError(w, http.StatusUnauthorized, ErrMissingToken)
			return
		}

		authResponse, err := verifyJWT(r.Header["Token"][0])
		if err != nil {
			httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
			return
		}

		if authResponse.Status == msgtypes.MISSING_TOKEN {
			httperrors.WriteJSONError(w, http.StatusUnauthorized, ErrMissingToken)
			return
		}

		if authResponse.Status == msgtypes.INVALID {
			httperrors.WriteJSONError(w, http.StatusUnauthorized, ErrInvalidToken)
			return
		}

		if authResponse.Status == msgtypes.NOT_FOUND {
			httperrors.WriteJSONError(w, http.StatusNotFound, ErrSubjectNotFound)
			return
		}

		r.Header.Add("Subject", strconv.Itoa(authResponse.ID))

		handler(w, r)
	})
}
