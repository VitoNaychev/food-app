package handlers_test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/testutil"
)

var restaurantHandlerMessage = "Hello from restaurant handler"
var addressHandlerMessage = "Hello from address handler"
var hoursHandlerMessage = "Hello from hours handler"
var menuHandlerMessage = "Hello from menu handler"

func fakeRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(restaurantHandlerMessage))
}

func fakeAddressHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(addressHandlerMessage))
}

func fakeHoursHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(hoursHandlerMessage))
}

func fakeMenuHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(menuHandlerMessage))
}

func TestRouterServer(t *testing.T) {
	fakeRestaurantHandler := http.HandlerFunc(fakeRestaurantHandler)
	fakeAddressServer := http.HandlerFunc(fakeAddressHandler)
	fakeHoursHandler := http.HandlerFunc(fakeHoursHandler)
	fakeMenuHandler := http.HandlerFunc(fakeMenuHandler)

	routerServer := handlers.NewRouterServer(
		fakeRestaurantHandler, fakeAddressServer, fakeHoursHandler, fakeMenuHandler)

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

	t.Run("routes requests to the hours server", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/restaurant/hours/", nil)
		response := httptest.NewRecorder()

		routerServer.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)

		want := hoursHandlerMessage
		got := getMessageFromBody(response.Body)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("routes requests to the menu server", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/restaurant/menu/", nil)
		response := httptest.NewRecorder()

		routerServer.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)

		want := menuHandlerMessage
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
