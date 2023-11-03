package integrationtest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/bt-customer-svc/handlers/address"
	"github.com/VitoNaychev/bt-customer-svc/handlers/customer"
	"github.com/VitoNaychev/bt-customer-svc/handlers/router"
	"github.com/VitoNaychev/bt-customer-svc/models"
	"github.com/VitoNaychev/bt-customer-svc/tests/testdata"
	"github.com/VitoNaychev/bt-customer-svc/tests/testutil"
)

func TestAddressServerOperations(t *testing.T) {
	connStr := SetupDatabaseContainer(t)

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	customerStore, err := models.NewPgCustomerStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	customerServer := customer.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, &customerStore)
	addressServer := address.NewCustomerAddressServer(&addressStore, &customerStore, testEnv.SecretKey)

	server := router.NewRouterServer(customerServer, addressServer)

	peterJWT := createNewCustomer(server, testdata.PeterCustomer)
	aliceJWT := createNewCustomer(server, testdata.AliceCustomer)

	createNewAddress(t, server, testdata.PeterAddress1, peterJWT)
	createNewAddress(t, server, testdata.PeterAddress2, peterJWT)
	createNewAddress(t, server, testdata.AliceAddress, aliceJWT)

	t.Run("get Peter's address", func(t *testing.T) {
		request := address.NewGetAddressRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := []address.GetAddressResponse{
			address.AddressToGetAddressResponse(testdata.PeterAddress1),
			address.AddressToGetAddressResponse(testdata.PeterAddress2),
		}
		var got []address.GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertAddresses(t, got, want)

	})

	t.Run("get Alice's address", func(t *testing.T) {
		request := address.NewGetAddressRequest(aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []address.GetAddressResponse{
			address.AddressToGetAddressResponse(testdata.AliceAddress),
		}
		var got []address.GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertAddresses(t, got, want)
	})

	t.Run("update Peter's address", func(t *testing.T) {
		updateAddress := testdata.PeterAddress2
		updateAddress.City = "Varna"

		request := address.NewUpdateAddressRequest(updateAddress, peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		request = address.NewGetAddressRequest(peterJWT)
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := []address.GetAddressResponse{
			address.AddressToGetAddressResponse(testdata.PeterAddress1),
			address.AddressToGetAddressResponse(updateAddress),
		}
		var got []address.GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertAddresses(t, got, want)
	})

	t.Run("delete Peter's address", func(t *testing.T) {
		deleteAddressId := testdata.PeterAddress2.Id

		request := address.NewDeleteAddressRequest(peterJWT, deleteAddressId)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		request = address.NewGetAddressRequest(peterJWT)
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := []address.GetAddressResponse{
			address.AddressToGetAddressResponse(testdata.PeterAddress1),
		}
		var got []address.GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertAddresses(t, got, want)
	})

}

func createNewAddress(t testing.TB, server http.Handler, a models.Address, customerJWT string) {
	request := address.NewCreateAddressRequest(customerJWT, a)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)
}
