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
		request := NewDeleteAddressRequest(invalidJWT, nil)
		request.Header.Add("Token", invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Bad Request on inavlid request", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)
		request := NewDeleteAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns Not Found on missing user", func(t *testing.T) {
		deleteAddressRequest := DeleteAddressRequest{Id: 0}
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(deleteAddressRequest)
		missingJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)

		request := NewDeleteAddressRequest(missingJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingCustomer)
	})

	t.Run("returns Not Found on missing address", func(t *testing.T) {
		deleteMissingAddressRequest := DeleteAddressRequest{Id: 10}
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(deleteMissingAddressRequest)
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewDeleteAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingAddress)
	})

	t.Run("returns Unathorized on delete on another customer's address", func(t *testing.T) {
		deleteAddressRequest := DeleteAddressRequest{td.AliceAddress.Id}
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(deleteAddressRequest)
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewDeleteAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})

	t.Run("deletes address on valid body and credentials", func(t *testing.T) {
		deleteAddressRequest := DeleteAddressRequest{td.PeterAddress1.Id}
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(deleteAddressRequest)
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewDeleteAddressRequest(peterJWT, body)
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
		request := NewAddAddressRequest(invalidJWT, nil)
		request.Header.Add("Token", invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Bad Request on inavlid request", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewAddAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns Not Found on missing user", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(td.AliceAddress)

		missingJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)

		request := NewAddAddressRequest(missingJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("saves Peter's new address", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(td.PeterAddress1)

		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := NewAddAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertStoredAddress(t, stubAddressStore, td.PeterAddress1)
	})

	t.Run("saves Alice's new address", func(t *testing.T) {
		stubAddressStore.Empty()

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(td.AliceAddress)

		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.AliceCustomer.Id)

		request := NewAddAddressRequest(peterJWT, body)
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
			addressToGetAddressResponse(td.PeterAddress1),
			addressToGetAddressResponse(td.PeterAddress2),
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
			addressToGetAddressResponse(td.AliceAddress),
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
