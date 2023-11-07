package unittest

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/handlers/auth"
	"github.com/VitoNaychev/bt-customer-svc/handlers/customer"
	"github.com/VitoNaychev/bt-customer-svc/handlers/validation"
	"github.com/VitoNaychev/bt-customer-svc/models"
	td "github.com/VitoNaychev/bt-customer-svc/tests/testdata"
	"github.com/VitoNaychev/bt-customer-svc/tests/testutil"
)

func TestUpdateUser(t *testing.T) {
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := customer.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("updates customer information on valid JWT", func(t *testing.T) {
		updateCustomer := td.PeterCustomer
		updateCustomer.FirstName = "John"
		updateCustomer.PhoneNumber = "+359 88 1234 213"

		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := customer.NewUpdateCustomerRequest(updateCustomer, peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertUpdatedCustomer(t, store, updateCustomer)
	})
}

func TestDeleteUser(t *testing.T) {
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := customer.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("deletes customer on valid JWT", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := customer.NewDeleteCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		testutil.AssertDeletedCustomer(t, store, td.PeterCustomer)
	})

	t.Run("returns Not Found on missing customer", func(t *testing.T) {
		missingCustomerJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)

		request := customer.NewDeleteCustomerRequest(missingCustomerJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingCustomer)
	})
}

func TestLoginUser(t *testing.T) {
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := customer.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("returns JWT on Peter's credentials", func(t *testing.T) {
		request := customer.NewLoginRequest(td.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertJWT(t, response.Header(), testEnv.SecretKey, td.PeterCustomer.Id)
	})

	t.Run("returns JWT on Alice's credentials", func(t *testing.T) {
		request := customer.NewLoginRequest(td.AliceCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertJWT(t, response.Header(), testEnv.SecretKey, td.AliceCustomer.Id)
	})

	t.Run("returns Unauthorized on invalid credentials", func(t *testing.T) {
		incorrectCustomer := td.PeterCustomer
		incorrectCustomer.Password = "passsword123"
		request := customer.NewLoginRequest(incorrectCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidCredentials)
	})

	t.Run("returns Unauthorized on missing user", func(t *testing.T) {
		missingCustomer := td.PeterCustomer
		missingCustomer.Email = "notanemail@gmail.com"
		request := customer.NewLoginRequest(missingCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingCustomer)
	})
}

func TestCreateUser(t *testing.T) {
	customerData := []models.Customer{}
	store := testutil.NewStubCustomerStore(customerData)
	server := customer.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("stores customer on POST", func(t *testing.T) {
		store.Empty()

		request := customer.NewCreateCustomerRequest(td.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusAccepted

		testutil.AssertStatus(t, got, want)
		testutil.AssertCreatedCustomer(t, store, td.PeterCustomer)
	})

	t.Run("returns JWT on POST", func(t *testing.T) {
		store.Empty()

		request := customer.NewCreateCustomerRequest(td.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusAccepted

		testutil.AssertStatus(t, got, want)
		testutil.AssertJWT(t, response.Header(), testEnv.SecretKey, td.PeterCustomer.Id)
	})

	t.Run("return Bad Request on user with same email", func(t *testing.T) {
		store.Empty()

		request := customer.NewCreateCustomerRequest(td.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		request = customer.NewCreateCustomerRequest(td.PeterCustomer)
		// reinit ResponseRecorder as it allows a
		// one-time only write of the Status Code
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrExistingUser)
	})
}

func TestGetUser(t *testing.T) {
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := customer.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("returns Peter's customer information", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)
		request := customer.NewGetCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[customer.GetCustomerResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		want := customer.CustomerToGetCustomerResponse(td.PeterCustomer)
		assertGetCustomerResponse(t, got, want)
	})

	t.Run("returns Alice's customer information", func(t *testing.T) {
		aliceJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.AliceCustomer.Id)
		request := customer.NewGetCustomerRequest(aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[customer.GetCustomerResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		want := customer.CustomerToGetCustomerResponse(td.AliceCustomer)
		assertGetCustomerResponse(t, got, want)
	})

	t.Run("returns Not Found on missing customer", func(t *testing.T) {
		noCustomerJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 3)
		request := customer.NewGetCustomerRequest(noCustomerJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingCustomer)
	})
}

func assertGetCustomerResponse(t testing.TB, got, want customer.GetCustomerResponse) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}