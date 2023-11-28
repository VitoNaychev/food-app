package testutil

import (
	"encoding/json"
	"io"
	"reflect"
	"strconv"
	"testing"

	"github.com/VitoNaychev/food-app/errorresponse"
	"github.com/golang-jwt/jwt/v5"
)

func AssertEqual[T any](t testing.TB, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertErrorResponse(t testing.TB, body io.Reader, expetedError error) {
	t.Helper()

	var errorResponse errorresponse.ErrorResponse
	json.NewDecoder(body).Decode(&errorResponse)

	if errorResponse.Message != expetedError.Error() {
		t.Errorf("got error %q want %q", errorResponse.Message, expetedError.Error())
	}
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %d, want %d", got, want)
	}
}

func AssertValidResponse(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("invalid response body %v", err)
	}
}

func AssertJWT(t testing.TB, jwtString string, secretKey []byte, wantId int) {
	t.Helper()

	token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {
		t.Fatalf("error verifying JWT: %v", err)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		t.Fatalf("did not get subject in JWT, expected %v", wantId)
	}

	gotId, err := strconv.Atoi(subject)
	if err != nil {
		t.Fatalf("did not get customer ID for subject, got %v", subject)
	}

	if gotId != wantId {
		t.Errorf("got customer id %v want %v", subject, wantId)
	}
}
