package integrationtest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/bt-customer-svc/handlers/customer"
	"github.com/VitoNaychev/bt-customer-svc/models"
	"github.com/VitoNaychev/bt-customer-svc/tests/testdata"
	"github.com/VitoNaychev/bt-customer-svc/tests/testutil"
)

func TestCustomerServerOperations(t *testing.T) {
	connStr := SetupDatabaseContainer(t)

	store, err := models.NewPgCustomerStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	server := customer.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, &store)

	peterJWT := createNewCustomer(server, testdata.PeterCustomer)
	aliceJWT := createNewCustomer(server, testdata.AliceCustomer)

	t.Run("retrieving customer", func(t *testing.T) {
		request := customer.NewGetCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := customer.CustomerToCustomerResponse(testdata.PeterCustomer)
		var got customer.CustomerResponse
		json.NewDecoder(response.Body).Decode(&got)
		testutil.AssertCustomerResponse(t, got, want)
	})

	t.Run("updating customer", func(t *testing.T) {
		updateCustomer := testdata.AliceCustomer
		updateCustomer.LastName = "Roper"

		request := customer.NewUpdateCustomerRequest(updateCustomer, aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := customer.CustomerToCustomerResponse(updateCustomer)
		var got customer.CustomerResponse
		json.NewDecoder(response.Body).Decode(&got)
		testutil.AssertCustomerResponse(t, got, want)
	})

	t.Run("deleting customer", func(t *testing.T) {
		request := customer.NewDeleteCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		request = customer.NewGetCustomerRequest(peterJWT)
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func createNewCustomer(server http.Handler, c models.Customer) string {
	request := customer.NewCreateCustomerRequest(c)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	return response.Header()["Token"][0]
}
