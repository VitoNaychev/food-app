package auth

import (
	"net/http"
	"strconv"
	"time"

	"github.com/VitoNaychev/food-app/errorresponse"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(secretKey []byte, expiresAt time.Duration, subject int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(int64(subject), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(jwtString string, secretKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {
		return nil, err
	}

	return token, nil
}

type Verifier interface {
	DoesSubjectExist(id int) (bool, error)
}

func AuthenticationMiddleware(endpointHandler func(w http.ResponseWriter, r *http.Request),
	verifier Verifier,
	secretKey []byte) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			errorresponse.WriteJSONError(w, http.StatusUnauthorized, ErrMissingToken)
			return
		}

		token, err := VerifyJWT(r.Header["Token"][0], secretKey)
		if err != nil {
			errorresponse.WriteJSONError(w, http.StatusUnauthorized, err)
			return
		}

		id, err := getIDFromToken(token)
		if err != nil {
			errorresponse.WriteJSONError(w, http.StatusUnauthorized, err)
			return
		}

		exists, err := verifier.DoesSubjectExist(id)
		if err != nil {
			errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
			return
		}

		if !exists {
			errorresponse.WriteJSONError(w, http.StatusNotFound, ErrSubjectNotFound)
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
