package auth

import (
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/golang-jwt/jwt/v5"
)

type Verifier interface {
	DoesSubjectExist(id int) (bool, error)
}

func AuthenticationMW(endpointHandler func(w http.ResponseWriter, r *http.Request),
	verifier Verifier,
	secretKey []byte) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			httperrors.WriteJSONError(w, http.StatusUnauthorized, ErrMissingToken)
			return
		}

		token, err := VerifyJWT(r.Header["Token"][0], secretKey)
		if err != nil {
			httperrors.WriteJSONError(w, http.StatusUnauthorized, err)
			return
		}

		id, err := getIDFromToken(token)
		if err != nil {
			httperrors.WriteJSONError(w, http.StatusUnauthorized, err)
			return
		}

		exists, err := verifier.DoesSubjectExist(id)
		if err != nil {
			httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
			return
		}

		if !exists {
			httperrors.WriteJSONError(w, http.StatusNotFound, ErrSubjectNotFound)
			return
		}

		r.Header.Add("Subject", strconv.Itoa(id))

		endpointHandler(w, r)
	})
}

func getIDFromToken(token *jwt.Token) (int, error) {
	subject, err := token.Claims.GetSubject()
	if err != nil || subject == "" {
		return -1, ErrMissingSubject
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return -1, ErrNonIntegerSubject
	}

	return id, nil
}
