package address

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/handlers/auth"
	as "github.com/VitoNaychev/bt-customer-svc/models/address_store"
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

func TestUpdateCustomerAddress(t *testing.T) {
	addressData := []as.Address{td.PeterAddress1, td.PeterAddress2, td.AliceAddress}
	customerData := []cs.Customer{td.PeterCustomer, td.AliceCustomer}
	stubAddressStore := testutil.NewStubAddressStore(addressData)
	stubCustomerStore := testutil.NewStubCustomerStore(customerData)
	server := CustomerAddressServer{stubAddressStore, stubCustomerStore, testEnv.SecretKey}

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		invalidJWT := "thisIsAnInvalidJWT"
		request, _ := http.NewRequest(http.MethodPut, "/customer/address/", nil)
		request.Header.Add("Token", invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("updates address on valid body and credentials", func(t *testing.T) {
		updatedAddress := td.PeterAddress2
		updatedAddress.City = "Varna"

		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewUpdateAddressRequest(updatedAddress, peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertUpdatedAddress(t, stubAddressStore, updatedAddress)
	})

	t.Run("returns Bad Request on invalid request", func(t *testing.T) {
		invalidAddress := as.Address{}

		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewUpdateAddressRequest(invalidAddress, peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidRequestField)
	})

	t.Run("returns Not Found on missing customer", func(t *testing.T) {
		updatedAddress := td.PeterAddress2
		updatedAddress.City = "Varna"

		missingJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)

		request := NewUpdateAddressRequest(updatedAddress, missingJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingCustomer)
	})

	t.Run("returns Not Found on missing address", func(t *testing.T) {
		updatedAddress := td.PeterAddress2
		updatedAddress.Id = 10

		missingJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewUpdateAddressRequest(updatedAddress, missingJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingAddress)
	})

	t.Run("returns Unauthorized on update on another customer's address", func(t *testing.T) {
		updatedAddress := td.PeterAddress2
		updatedAddress.City = "Varna"

		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.AliceCustomer.Id)

		request := NewUpdateAddressRequest(updatedAddress, peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})
}

func TestDeleteCustomerAddress(t *testing.T) {
	addressData := []as.Address{td.PeterAddress1, td.PeterAddress2, td.AliceAddress}
	customerData := []cs.Customer{td.PeterCustomer, td.AliceCustomer}
	stubAddressStore := testutil.NewStubAddressStore(addressData)
	stubCustomerStore := testutil.NewStubCustomerStore(customerData)
	server := CustomerAddressServer{stubAddressStore, stubCustomerStore, testEnv.SecretKey}

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		invalidJWT := "thisIsAnInvalidJWT"
		request := NewDeleteAddressRequest(invalidJWT, td.PeterAddress1.Id)
		request.Header.Add("Token", invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Bad Request on inavlid request", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request, _ := http.NewRequest(http.MethodDelete, "/customer/address", body)
		request.Header.Add("Token", peterJWT)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns Not Found on missing user", func(t *testing.T) {
		missingJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)

		request := NewDeleteAddressRequest(missingJWT, td.PeterAddress1.Id)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingCustomer)
	})

	t.Run("returns Not Found on missing address", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewDeleteAddressRequest(peterJWT, 10)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingAddress)
	})

	t.Run("returns Unathorized on delete on another customer's address", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewDeleteAddressRequest(peterJWT, td.AliceAddress.Id)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})

	t.Run("deletes address on valid body and credentials", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewDeleteAddressRequest(peterJWT, td.PeterAddress1.Id)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertDeletedAddress(t, stubAddressStore, td.PeterAddress1)
	})
}

func TestSaveCustomerAddress(t *testing.T) {
	addressData := []as.Address{}
	customerData := []cs.Customer{td.PeterCustomer, td.AliceCustomer}
	stubAddressStore := testutil.NewStubAddressStore(addressData)
	stubCustomerStore := testutil.NewStubCustomerStore(customerData)
	server := CustomerAddressServer{stubAddressStore, stubCustomerStore, testEnv.SecretKey}

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		invalidJWT := "thisIsAnInvalidJWT"
		request := NewAddAddressRequest(invalidJWT, td.AliceAddress)
		request.Header.Add("Token", invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Bad Request on inavlid request", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewAddAddressRequest(peterJWT, as.Address{})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns Not Found on missing user", func(t *testing.T) {
		missingJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)

		request := NewAddAddressRequest(missingJWT, td.AliceAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("saves Peter's new address", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewAddAddressRequest(peterJWT, td.PeterAddress1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertStoredAddress(t, stubAddressStore, td.PeterAddress1)
	})

	t.Run("saves Alice's new address", func(t *testing.T) {
		stubAddressStore.Empty()

		aliceJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.AliceCustomer.Id)

		request := NewAddAddressRequest(aliceJWT, td.AliceAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertStoredAddress(t, stubAddressStore, td.AliceAddress)
	})
}

func TestGetCustomerAddress(t *testing.T) {
	addressData := []as.Address{td.PeterAddress1, td.PeterAddress2, td.AliceAddress}
	customerData := []cs.Customer{td.PeterCustomer, td.AliceCustomer}
	stubAddressStore := testutil.NewStubAddressStore(addressData)
	stubCustomerStore := testutil.NewStubCustomerStore(customerData)
	server := CustomerAddressServer{stubAddressStore, stubCustomerStore, testEnv.SecretKey}

	t.Run("returns Peter's addresses", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)
		request := NewGetAddressRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []GetAddressResponse{
			AddressToGetAddressResponse(td.PeterAddress1),
			AddressToGetAddressResponse(td.PeterAddress2),
		}
		var got []GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)

		assertAddresses(t, got, want)
	})

	t.Run("returns Alice's addresses", func(t *testing.T) {
		aliceJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.AliceCustomer.Id)
		request := NewGetAddressRequest(aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []GetAddressResponse{
			AddressToGetAddressResponse(td.AliceAddress),
		}
		var got []GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)

		assertAddresses(t, got, want)
	})

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		invalidJWT := "thisIsAnInvalidJWT"
		request := NewGetAddressRequest(invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Not Found on missing user", func(t *testing.T) {
		aliceJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)
		request := NewGetAddressRequest(aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingCustomer)
	})
}

func assertAddresses(t testing.TB, got, want []GetAddressResponse) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, td.AliceAddress)
	}
}
