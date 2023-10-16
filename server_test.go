package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestCreateUser(t *testing.T) {
	t.Run("stores customer on POST and returns JWT", func(t *testing.T) {
		customer := CreateCustomerRequest{
			FirstName:   "Peter",
			LastName:    "Smith",
			PhoneNumber: "+359 88 576 5981",
			Email:       "petesmith@gmail.com",
			Password:    "samoMBTbratmeeee",
		}

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(customer)

		request, _ := http.NewRequest(http.MethodPost, "/customer/", body)
		response := httptest.NewRecorder()

		CustomerServer(response, request)

		got := response.Code
		want := http.StatusAccepted

		assertStatus(t, got, want)

	})
}

// func newCreateUserRequest()

func TestGetUser(t *testing.T) {
	secretKey := []byte("mySecretKey")

	peterJWT, _ := generateJWT(secretKey, 0)
	aliceJWT, _ := generateJWT(secretKey, 1)

	t.Run("returns Peter's customer information", func(t *testing.T) {
		request := newGetCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		CustomerServer(response, request)

		var got GetCustomerResponse
		json.NewDecoder(response.Body).Decode(&got)

		want := GetCustomerResponse{
			FirstName:   "Peter",
			LastName:    "Smith",
			PhoneNumber: "+359 88 576 5981",
			Email:       "petesmith@gmail.com",
		}

		assertGetCustomerResponse(t, got, want)
	})

	t.Run("returns Alice's customer information", func(t *testing.T) {
		request := newGetCustomerRequest(aliceJWT)
		response := httptest.NewRecorder()

		CustomerServer(response, request)

		var got GetCustomerResponse
		json.NewDecoder(response.Body).Decode(&got)

		want := GetCustomerResponse{
			FirstName:   "Alice",
			LastName:    "Johnson",
			PhoneNumber: "+359 88 444 2222",
			Email:       "alicejohn@gmail.com",
		}

		assertGetCustomerResponse(t, got, want)
	})
}

func newGetCustomerRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/customer/", nil)
	request.Header.Add("Token", jwt)

	return request
}

func generateJWT(secretKey []byte, subject int) (string, error) {
	tenSecondsFronNow := time.Now().Add(10 * time.Second)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(int64(subject), 10),
		ExpiresAt: jwt.NewNumericDate(tenSecondsFronNow),
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}

func assertGetCustomerResponse(t testing.TB, got, want GetCustomerResponse) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
