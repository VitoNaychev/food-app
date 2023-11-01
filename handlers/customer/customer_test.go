package customer

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/handlers/auth"
	"github.com/VitoNaychev/bt-customer-svc/handlers/validation"
	cs "github.com/VitoNaychev/bt-customer-svc/models/customer_store"
	td "github.com/VitoNaychev/bt-customer-svc/testdata"
	"github.com/VitoNaychev/bt-customer-svc/testutil"
)

var testEnv handlers.TestEnv

func TestMain(m *testing.M) {
	testEnv = handlers.LoadTestEnviornment()

	code := m.Run()
	os.Exit(code)
}

func TestUpdateUser(t *testing.T) {
	customerData := []cs.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("updates customer information on valid JWT", func(t *testing.T) {
		customer := td.PeterCustomer
		customer.FirstName = "John"
		customer.PhoneNumber = "+359 88 1234 213"

		request := NewUpdateCustomerRequest(customer, testEnv.SecretKey, testEnv.ExpiresAt)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertUpdatedCustomer(t, store, customer)
	})
}

func TestDeleteUser(t *testing.T) {
	customerData := []cs.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("deletes customer on valid JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/customer/", nil)
		response := httptest.NewRecorder()

		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)
		request.Header.Add("Token", peterJWT)

		server.ServeHTTP(response, request)
		testutil.AssertDeletedCustomer(t, store, td.PeterCustomer)
	})

	t.Run("returns Not Found on missing customer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/customer/", nil)
		response := httptest.NewRecorder()

		missingCustomerJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)
		request.Header.Add("Token", missingCustomerJWT)

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingCustomer)
	})
}

func TestLoginUser(t *testing.T) {
	customerData := []cs.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("returns JWT on Peter's credentials", func(t *testing.T) {
		request := NewLoginRequest(td.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertJWT(t, response.Header(), testEnv.SecretKey, td.PeterCustomer.Id)
	})

	t.Run("returns JWT on Alice's credentials", func(t *testing.T) {
		request := NewLoginRequest(td.AliceCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertJWT(t, response.Header(), testEnv.SecretKey, td.AliceCustomer.Id)
	})

	t.Run("returns Unauthorized on invalid credentials", func(t *testing.T) {
		incorrectCustomer := td.PeterCustomer
		incorrectCustomer.Password = "passsword123"
		request := NewLoginRequest(incorrectCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidCredentials)
	})

	t.Run("returns Unauthorized on missing user", func(t *testing.T) {
		missingCustomer := td.PeterCustomer
		missingCustomer.Email = "notanemail@gmail.com"
		request := NewLoginRequest(missingCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingCustomer)
	})
}

func TestCreateUser(t *testing.T) {
	customerData := []cs.Customer{}
	store := testutil.NewStubCustomerStore(customerData)
	server := NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("stores customer on POST", func(t *testing.T) {
		store.Empty()

		request := NewCreateCustomerRequest(td.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusAccepted

		testutil.AssertStatus(t, got, want)
		testutil.AssertStoredCustomer(t, store, td.PeterCustomer)
	})

	t.Run("returns JWT on POST", func(t *testing.T) {
		store.Empty()

		request := NewCreateCustomerRequest(td.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusAccepted

		testutil.AssertStatus(t, got, want)
		testutil.AssertJWT(t, response.Header(), testEnv.SecretKey, td.PeterCustomer.Id)
	})

	t.Run("return Bad Request on user with same email", func(t *testing.T) {
		store.Empty()

		request := NewCreateCustomerRequest(td.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		request = NewCreateCustomerRequest(td.PeterCustomer)
		// reinit ResponseRecorder as it allows a
		// one-time only write of the Status Code
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrExistingUser)
	})
}

func TestGetUser(t *testing.T) {
	customerData := []cs.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("returns Peter's customer information", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)
		request := NewGetCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		var got GetCustomerResponse
		err := validation.ValidateBody(response.Body, &got)
		testutil.AssertValidResponse(t, err)

		want := customerToGetCustomerResponse(td.PeterCustomer)
		assertGetCustomerResponse(t, got, want)
	})

	t.Run("returns Alice's customer information", func(t *testing.T) {
		aliceJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.AliceCustomer.Id)
		request := NewGetCustomerRequest(aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		var got GetCustomerResponse
		err := validation.ValidateBody(response.Body, &got)
		testutil.AssertValidResponse(t, err)

		want := customerToGetCustomerResponse(td.AliceCustomer)
		assertGetCustomerResponse(t, got, want)
	})

	t.Run("returns Not Found on missing customer", func(t *testing.T) {
		noCustomerJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 3)
		request := NewGetCustomerRequest(noCustomerJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingCustomer)
	})
}

func assertGetCustomerResponse(t testing.TB, got, want GetCustomerResponse) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
