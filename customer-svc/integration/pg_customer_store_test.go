package integrationtest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/customer-svc/handlers"
	"github.com/VitoNaychev/food-app/customer-svc/models"
	"github.com/VitoNaychev/food-app/customer-svc/testdata"
	"github.com/VitoNaychev/food-app/customer-svc/testutil"
)

func TestCustomerServerOperations(t *testing.T) {
	connStr := SetupDatabaseContainer(t)

	store, err := models.NewPgCustomerStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	server := handlers.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, &store)

	var peterJWT string
	var createdSuccessfully bool

	createdSuccessfully = t.Run("create new customer", func(t *testing.T) {
		request := handlers.NewCreateCustomerRequest(testdata.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)

		wantCustomer := handlers.CustomerToCustomerResponse(testdata.PeterCustomer)
		got := testutil.ParseCreateCustomerResponse(t, response.Body)

		testutil.AssertEqual(t, got.Customer, wantCustomer)

		peterJWT = got.JWT.Token
	})

	if createdSuccessfully {
		t.Run("retrieve customer", func(t *testing.T) {
			request := handlers.NewGetCustomerRequest(peterJWT)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			testutil.AssertStatus(t, response.Code, http.StatusOK)

			want := handlers.CustomerToCustomerResponse(testdata.PeterCustomer)
			got := testutil.ParseCustomerResponse(t, response.Body)

			testutil.AssertEqual(t, got, want)
		})

		t.Run("update customer", func(t *testing.T) {
			updateCustomer := testdata.PeterCustomer
			updateCustomer.LastName = "Roper"
			updateCustomer.Email = "peteroper@gmail.com"

			request := handlers.NewUpdateCustomerRequest(updateCustomer, peterJWT)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			testutil.AssertStatus(t, response.Code, http.StatusOK)

			want := handlers.CustomerToCustomerResponse(updateCustomer)
			got := testutil.ParseCustomerResponse(t, response.Body)

			testutil.AssertEqual(t, got, want)
		})

		t.Run("delete customer", func(t *testing.T) {
			request := handlers.NewDeleteCustomerRequest(peterJWT)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			testutil.AssertStatus(t, response.Code, http.StatusOK)

		})

		t.Run("retrieve deleted customer", func(t *testing.T) {
			request := handlers.NewGetCustomerRequest(peterJWT)
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			testutil.AssertStatus(t, response.Code, http.StatusNotFound)
			testutil.AssertErrorResponse(t, response.Body, handlers.ErrCustomerNotFound)
		})
	}
}
