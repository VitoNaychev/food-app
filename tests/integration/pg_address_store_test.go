package integrationtest

import (
	"context"
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

	var createdSuccessfully bool

	createdSuccessfully = t.Run("create addresses", func(t *testing.T) {
		response := createNewAddress(t, server, testdata.PeterAddress1, peterJWT)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got := testutil.ParseAddressResponse(t, response.Body)
		testutil.AssertAddressResponse(t, got, testdata.PeterAddress1)

		response = createNewAddress(t, server, testdata.PeterAddress2, peterJWT)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got = testutil.ParseAddressResponse(t, response.Body)
		testutil.AssertAddressResponse(t, got, testdata.PeterAddress2)

	})

	if !createdSuccessfully {
		return
	}

	t.Run("get addresses", func(t *testing.T) {
		request := address.NewGetAddressRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := []address.GetAddressResponse{
			address.AddressToGetAddressResponse(testdata.PeterAddress1),
			address.AddressToGetAddressResponse(testdata.PeterAddress2),
		}
		got := testutil.ParseGetAddressResponse(t, response.Body)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertGetAddressResponse(t, got, want)
	})

	t.Run("update address", func(t *testing.T) {
		updateAddress := testdata.PeterAddress2
		updateAddress.City = "Varna"

		request := address.NewUpdateAddressRequest(peterJWT, updateAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got := testutil.ParseAddressResponse(t, response.Body)
		testutil.AssertAddressResponse(t, got, updateAddress)
	})

	t.Run("delete address", func(t *testing.T) {
		deleteAddressRequest := address.DeleteAddressRequest{Id: testdata.PeterAddress2.Id}

		request := address.NewDeleteAddressRequest(peterJWT, deleteAddressRequest)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("get reamaining address", func(t *testing.T) {
		request := address.NewGetAddressRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := []address.GetAddressResponse{
			address.AddressToGetAddressResponse(testdata.PeterAddress1),
		}
		got := testutil.ParseGetAddressResponse(t, response.Body)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertGetAddressResponse(t, got, want)
	})

}

func createNewCustomer(server http.Handler, c models.Customer) string {
	request := customer.NewCreateCustomerRequest(c)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	return response.Header()["Token"][0]
}

func createNewAddress(t testing.TB, server http.Handler, a models.Address, customerJWT string) *httptest.ResponseRecorder {
	request := address.NewCreateAddressRequest(customerJWT, a)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	return response
}
