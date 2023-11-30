package tabletests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/order-svc/testutil"
)

func RunAuthenticationTests(t *testing.T, server http.Handler, cases map[string]*http.Request) {
	t.Helper()

	invalidJWT := "thisIsAnInvalidJWT"
	for name, request := range cases {
		t.Run(name, func(t *testing.T) {
			t.Helper()

			request.Header.Add("Token", invalidJWT)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		})
	}
}

type DummyRequest struct {
	S string
}

func NewDummyRequest(method string, path string, jwt string) *http.Request {
	dummyRequest := DummyRequest{"validation test"}

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(dummyRequest)

	request, _ := http.NewRequest(method, path, body)
	request.Header.Add("Token", jwt)

	return request
}

func RunRequestValidationTests(t *testing.T, server http.Handler, cases map[string]*http.Request) {
	t.Helper()

	for name, request := range cases {
		t.Run(name, func(t *testing.T) {
			t.Helper()

			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		})
	}
}

type GenericResponse interface{}

type ResponseValidationTestcase struct {
	Name               string
	Request            *http.Request
	ValidationFunction func(io.Reader) (GenericResponse, error)
}

func RunResponseValidationTests(t *testing.T, server http.Handler, cases []ResponseValidationTestcase) {
	t.Helper()

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			t.Helper()

			response := httptest.NewRecorder()

			server.ServeHTTP(response, test.Request)

			_, err := test.ValidationFunction(response.Body)
			if err != nil {
				t.Errorf("invalid response body, %v", err)
			}
		})
	}
}
