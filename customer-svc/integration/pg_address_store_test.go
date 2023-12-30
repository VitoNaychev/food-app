package integrationtest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/customer-svc/handlers"
	"github.com/VitoNaychev/food-app/customer-svc/models"
	"github.com/VitoNaychev/food-app/customer-svc/testdata"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/parser"
	"github.com/VitoNaychev/food-app/pgconfig"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestAddressServerOperations(t *testing.T) {
	config := pgconfig.GetConfigFromEnv(testEnv)
	integrationutil.SetupDatabaseContainer(t, &config)

	connStr := config.GetConnectionString()

	addressStore, err := models.NewPgAddressStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	customerStore, err := models.NewPgCustomerStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	customerServer := handlers.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, &customerStore)
	addressServer := handlers.NewCustomerAddressServer(&addressStore, &customerStore, testEnv.SecretKey)

	server := handlers.NewRouterServer(customerServer, addressServer)

	peterJWT := createNewCustomer(server, testdata.PeterCustomer)

	var createdSuccessfully bool

	createdSuccessfully = t.Run("create addresses", func(t *testing.T) {
		response := createNewAddress(t, server, testdata.PeterAddress1, peterJWT)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got := parser.FromJSON[models.Address](response.Body)
		testutil.AssertEqual(t, got, testdata.PeterAddress1)

		response = createNewAddress(t, server, testdata.PeterAddress2, peterJWT)
		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got = parser.FromJSON[models.Address](response.Body)
		testutil.AssertEqual(t, got, testdata.PeterAddress2)

	})

	if !createdSuccessfully {
		return
	}

	t.Run("get addresses", func(t *testing.T) {
		request := handlers.NewGetAddressRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := []handlers.GetAddressResponse{
			handlers.AddressToGetAddressResponse(testdata.PeterAddress1),
			handlers.AddressToGetAddressResponse(testdata.PeterAddress2),
		}
		got := parser.FromJSON[[]handlers.GetAddressResponse](response.Body)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertEqual(t, got, want)
	})

	t.Run("update address", func(t *testing.T) {
		updateAddress := testdata.PeterAddress2
		updateAddress.City = "Varna"

		request := handlers.NewUpdateAddressRequest(peterJWT, updateAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got := parser.FromJSON[models.Address](response.Body)
		testutil.AssertEqual(t, got, updateAddress)
	})

	t.Run("delete address", func(t *testing.T) {
		deleteAddressRequest := handlers.DeleteAddressRequest{Id: testdata.PeterAddress2.Id}

		request := handlers.NewDeleteAddressRequest(peterJWT, deleteAddressRequest)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
	})

	t.Run("get reamaining address", func(t *testing.T) {
		request := handlers.NewGetAddressRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := []handlers.GetAddressResponse{
			handlers.AddressToGetAddressResponse(testdata.PeterAddress1),
		}
		got := parser.FromJSON[[]handlers.GetAddressResponse](response.Body)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertEqual(t, got, want)
	})

}

func createNewCustomer(server http.Handler, c models.Customer) string {
	request := handlers.NewCreateCustomerRequest(c)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	var createCustomerResponse handlers.CreateCustomerResponse
	json.NewDecoder(response.Body).Decode(&createCustomerResponse)
	return createCustomerResponse.JWT.Token
}

func createNewAddress(t testing.TB, server http.Handler, a models.Address, customerJWT string) *httptest.ResponseRecorder {
	request := handlers.NewCreateAddressRequest(customerJWT, a)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	return response
}
