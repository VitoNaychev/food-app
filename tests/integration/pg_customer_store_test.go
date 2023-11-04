package integrationtest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
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

	var peterJWT string
	var createdSuccessfully bool

	createdSuccessfully = t.Run("create new customer", func(t *testing.T) {
		request := customer.NewCreateCustomerRequest(testdata.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)

		want := customer.CustomerToCustomerResponse(testdata.PeterCustomer)
		got := testutil.ParseCustomerResponse(t, response.Body)

		testutil.AssertCustomerResponse(t, got, want)

		peterJWT = response.Header()["Token"][0]
	})

	if !createdSuccessfully {
		return
	}

	t.Run("retrieve customer", func(t *testing.T) {
		request := customer.NewGetCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := customer.CustomerToCustomerResponse(testdata.PeterCustomer)
		got := testutil.ParseCustomerResponse(t, response.Body)

		testutil.AssertCustomerResponse(t, got, want)
	})

	t.Run("update customer", func(t *testing.T) {
		updateCustomer := testdata.PeterCustomer
		updateCustomer.LastName = "Roper"
		updateCustomer.Email = "peteroper@gmail.com"

		request := customer.NewUpdateCustomerRequest(updateCustomer, peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := customer.CustomerToCustomerResponse(updateCustomer)
		got := testutil.ParseCustomerResponse(t, response.Body)

		testutil.AssertCustomerResponse(t, got, want)
	})

	t.Run("delete customer", func(t *testing.T) {
		request := customer.NewDeleteCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

	})

	t.Run("retrieve deleted customer", func(t *testing.T) {
		request := customer.NewGetCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingCustomer)
	})
}
