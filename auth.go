package main

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingToken = errors.New("missing token")
)

func generateJWT(secretKey []byte, expiresAt time.Time, subject int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(int64(subject), 10),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyJWT(header http.Header, secretKey []byte) (*jwt.Token, error) {
	if header["Token"] == nil {
		return nil, ErrMissingToken
	}

	token, err := jwt.Parse(header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {
		return nil, err
	}

	return token, nil
}
