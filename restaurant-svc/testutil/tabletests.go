package testutil

import (
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
