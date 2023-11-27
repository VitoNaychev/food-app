package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/testutil"
)

var restaurantHandlerMessage = "Hello from restaurant handler"
var addressHandlerMessage = "Hello from address handler"

func fakeRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(restaurantHandlerMessage))
}

func fakeAddressHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(addressHandlerMessage))

}

func TestRouterServer(t *testing.T) {
	fakeRestaurantHandler := http.HandlerFunc(fakeRestaurantHandler)
	fakeAddressServer := http.HandlerFunc(fakeAddressHandler)

	routerServer := handlers.NewRouterServer(fakeRestaurantHandler, fakeAddressServer)

	t.Run("routes requests to the restaurant server", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/restaurant/", nil)
		response := httptest.NewRecorder()

		routerServer.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)

		want := restaurantHandlerMessage
		got := getMessageFromBody(response.Body)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("routes requests to the address server", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/restaurant/address/", nil)
		response := httptest.NewRecorder()

		routerServer.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)

		want := addressHandlerMessage
		got := getMessageFromBody(response.Body)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns Not Found on unknown path", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/unknown/path", nil)
		response := httptest.NewRecorder()

		routerServer.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func getMessageFromBody(body io.Reader) string {
	buf := bytes.NewBuffer([]byte{})
	buf.ReadFrom(body)
	return buf.String()
}
