package unittest

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/bt-customer-svc/handlers/router"
	"github.com/VitoNaychev/bt-customer-svc/tests/testutil"
)

var customerHandlerMessage = "Hello from customer handler"
var addressHandlerMessage = "Hello from address handler"

func fakeCustomerHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(customerHandlerMessage))
}

func fakeAddressHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(addressHandlerMessage))

}

func TestRouterServer(t *testing.T) {
	fakeCustomerServer := http.HandlerFunc(fakeCustomerHandler)
	fakeAddressServer := http.HandlerFunc(fakeAddressHandler)

	routerServer := router.InitRouterServer(fakeCustomerServer, fakeAddressServer)

	t.Run("routes requests to the customer server", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/customer/", nil)
		response := httptest.NewRecorder()

		routerServer.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)

		want := customerHandlerMessage
		got := getMessageFromBody(response.Body)

		assertHandlerMessage(t, got, want)
	})

	t.Run("routes requests to the address server", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/customer/address/", nil)
		response := httptest.NewRecorder()

		routerServer.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)

		want := addressHandlerMessage
		got := getMessageFromBody(response.Body)

		assertHandlerMessage(t, got, want)
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

func assertHandlerMessage(t testing.TB, got, want string) {
	t.Helper()

	if got != want {
		t.Errorf("got message %q want %q", got, want)
	}
}
