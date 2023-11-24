package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/customer-svc/handlers"
	"github.com/VitoNaychev/food-app/customer-svc/models"
	td "github.com/VitoNaychev/food-app/customer-svc/testdata"
	"github.com/VitoNaychev/food-app/customer-svc/testutil"
)

func TestAddressEndpointAuthentication(t *testing.T) {
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	stubAddressStore := testutil.NewStubAddressStore(nil)
	stubCustomerStore := testutil.NewStubCustomerStore(customerData)
	server := handlers.NewCustomerAddressServer(stubAddressStore, stubCustomerStore, testEnv.SecretKey)

	invalidJWT := "thisIsAnInvalidJWT"
	cases := map[string]*http.Request{
		"get address authentication":    handlers.NewGetAddressRequest(invalidJWT),
		"create address authentication": handlers.NewCreateAddressRequest(invalidJWT, models.Address{}),
		"update address authentication": handlers.NewUpdateAddressRequest(invalidJWT, models.Address{}),
		"delete address authentication": handlers.NewDeleteAddressRequest(invalidJWT, handlers.DeleteAddressRequest{}),
	}

	for name, request := range cases {
		t.Run(name, func(t *testing.T) {
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		})
	}
}

func TestUpdateCustomerAddress(t *testing.T) {
	addressData := []models.Address{td.PeterAddress1, td.PeterAddress2, td.AliceAddress}
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	stubAddressStore := testutil.NewStubAddressStore(addressData)
	stubCustomerStore := testutil.NewStubCustomerStore(customerData)
	server := handlers.NewCustomerAddressServer(stubAddressStore, stubCustomerStore, testEnv.SecretKey)

	t.Run("updates address on valid body and credentials", func(t *testing.T) {
		updatedAddress := td.PeterAddress2
		updatedAddress.City = "Varna"

		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := handlers.NewUpdateAddressRequest(peterJWT, updatedAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertUpdatedAddress(t, stubAddressStore, updatedAddress)
	})

	t.Run("returns Bad Request on invalid request", func(t *testing.T) {
		invalidAddress := models.Address{}

		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := handlers.NewUpdateAddressRequest(peterJWT, invalidAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidRequestField)
	})

	t.Run("returns Not Found on missing customer", func(t *testing.T) {
		updatedAddress := td.PeterAddress2
		updatedAddress.City = "Varna"

		missingJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)

		request := handlers.NewUpdateAddressRequest(missingJWT, updatedAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrCustomerNotFound)
	})

	t.Run("returns Not Found on missing address", func(t *testing.T) {
		updatedAddress := td.PeterAddress2
		updatedAddress.Id = 10

		missingJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := handlers.NewUpdateAddressRequest(missingJWT, updatedAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingAddress)
	})

	t.Run("returns Unauthorized on update on another customer's address", func(t *testing.T) {
		updatedAddress := td.PeterAddress2
		updatedAddress.City = "Varna"

		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.AliceCustomer.Id)

		request := handlers.NewUpdateAddressRequest(peterJWT, updatedAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})
}

func TestDeleteCustomerAddress(t *testing.T) {
	addressData := []models.Address{td.PeterAddress1, td.PeterAddress2, td.AliceAddress}
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	stubAddressStore := testutil.NewStubAddressStore(addressData)
	stubCustomerStore := testutil.NewStubCustomerStore(customerData)
	server := handlers.NewCustomerAddressServer(stubAddressStore, stubCustomerStore, testEnv.SecretKey)

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
		deleteAddressRequest := handlers.DeleteAddressRequest{Id: td.PeterAddress1.Id}

		request := handlers.NewDeleteAddressRequest(missingJWT, deleteAddressRequest)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrCustomerNotFound)
	})

	t.Run("returns Not Found on missing address", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)
		deleteAddressRequest := handlers.DeleteAddressRequest{Id: 10}

		request := handlers.NewDeleteAddressRequest(peterJWT, deleteAddressRequest)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingAddress)
	})

	t.Run("returns Unathorized on delete on another customer's address", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)
		deleteAddressRequest := handlers.DeleteAddressRequest{Id: td.AliceAddress.Id}

		request := handlers.NewDeleteAddressRequest(peterJWT, deleteAddressRequest)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})

	t.Run("deletes address on valid body and credentials", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)
		deleteAddressRequest := handlers.DeleteAddressRequest{Id: td.PeterAddress1.Id}

		request := handlers.NewDeleteAddressRequest(peterJWT, deleteAddressRequest)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertDeletedAddress(t, stubAddressStore, td.PeterAddress1)
	})
}

func TestSaveCustomerAddress(t *testing.T) {
	addressData := []models.Address{}
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	stubAddressStore := testutil.NewStubAddressStore(addressData)
	stubCustomerStore := testutil.NewStubCustomerStore(customerData)
	server := handlers.NewCustomerAddressServer(stubAddressStore, stubCustomerStore, testEnv.SecretKey)

	t.Run("returns Bad Request on inavlid request", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := handlers.NewCreateAddressRequest(peterJWT, models.Address{})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns Not Found on missing user", func(t *testing.T) {
		missingJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)

		request := handlers.NewCreateAddressRequest(missingJWT, td.AliceAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("saves Peter's new address", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := handlers.NewCreateAddressRequest(peterJWT, td.PeterAddress1)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertStoredAddress(t, stubAddressStore, td.PeterAddress1)
	})

	t.Run("saves Alice's new address", func(t *testing.T) {
		stubAddressStore.Empty()

		aliceJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.AliceCustomer.Id)

		request := handlers.NewCreateAddressRequest(aliceJWT, td.AliceAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertStoredAddress(t, stubAddressStore, td.AliceAddress)
	})
}

func TestGetCustomerAddress(t *testing.T) {
	addressData := []models.Address{td.PeterAddress1, td.PeterAddress2, td.AliceAddress}
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	stubAddressStore := testutil.NewStubAddressStore(addressData)
	stubCustomerStore := testutil.NewStubCustomerStore(customerData)
	server := handlers.NewCustomerAddressServer(stubAddressStore, stubCustomerStore, testEnv.SecretKey)

	t.Run("returns Peter's addresses", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)
		request := handlers.NewGetAddressRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []handlers.GetAddressResponse{
			handlers.AddressToGetAddressResponse(td.PeterAddress1),
			handlers.AddressToGetAddressResponse(td.PeterAddress2),
		}
		var got []handlers.GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns Alice's addresses", func(t *testing.T) {
		aliceJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.AliceCustomer.Id)
		request := handlers.NewGetAddressRequest(aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []handlers.GetAddressResponse{
			handlers.AddressToGetAddressResponse(td.AliceAddress),
		}
		var got []handlers.GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns Not Found on missing user", func(t *testing.T) {
		aliceJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)
		request := handlers.NewGetAddressRequest(aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrCustomerNotFound)
	})
}
