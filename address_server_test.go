package bt_customer_svc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

type StubAddressStore struct {
	addresses        []Address
	savedAddress     Address
	deletedAddressId int
}

func (s *StubAddressStore) GetAddressesByCustomerId(customerId int) ([]Address, error) {
	if customerId == peterCustomer.Id {
		return []Address{peterAddress1, peterAddress2}, nil
	}

	if customerId == aliceCustomer.Id {
		return []Address{aliceAddress}, nil
	}

	return []Address{}, nil
}

func (s *StubAddressStore) StoreAddress(address Address) {
	s.savedAddress = address
}

func (s *StubAddressStore) DeleteAddressById(id int) error {
	_, err := s.GetAddressById(id)
	if err != nil {
		return err
	} else {
		s.deletedAddressId = id
		return nil
	}
}

func (s *StubAddressStore) GetAddressById(id int) (Address, error) {
	for _, address := range s.addresses {
		if address.Id == id {
			return address, nil
		}
	}
	return Address{}, fmt.Errorf("address with id %d doesn't exist", id)
}

func TestDeleteCustomerAddress(t *testing.T) {
	stubAddressStore := &StubAddressStore{
		[]Address{peterAddress1, peterAddress2, aliceAddress}, Address{}, 0,
	}
	stubCustomerStore := &StubCustomerStore{[]Customer{peterCustomer, aliceCustomer}, nil, nil}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	expiresAt := time.Now().Add(time.Second)
	server := CustomerAddressServer{stubAddressStore, stubCustomerStore, secretKey}

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		invalidJWT := "thisIsAnInvalidJWT"
		request := newDeleteAddressRequest(invalidJWT, nil)
		request.Header.Add("Token", invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Bad Request on inavlid request", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, peterCustomer.Id)
		request := newDeleteAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns Not Found on missing user", func(t *testing.T) {
		deleteAddressRequest := DeleteAddressRequest{Id: 0}
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(deleteAddressRequest)
		missingJWT, _ := GenerateJWT(secretKey, expiresAt, 10)

		request := newDeleteAddressRequest(missingJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		assertErrorResponse(t, errorResponse, ErrMissingCustomer)
	})

	t.Run("returns Not Found on missing address", func(t *testing.T) {
		deleteMissingAddressRequest := DeleteAddressRequest{Id: 10}
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(deleteMissingAddressRequest)
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, peterCustomer.Id)

		request := newDeleteAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		assertErrorResponse(t, errorResponse, ErrMissingAddress)
	})

	t.Run("returns Unathorized on attempt to delete another customer's address", func(t *testing.T) {
		deleteAddressRequest := DeleteAddressRequest{Id: 2}
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(deleteAddressRequest)
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, peterCustomer.Id)

		request := newDeleteAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		assertErrorResponse(t, errorResponse, ErrUnathorizedAction)
	})

	t.Run("deletes address on valid body and credentials", func(t *testing.T) {
		deleteAddressRequest := DeleteAddressRequest{Id: 1}
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(deleteAddressRequest)
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, peterCustomer.Id)

		request := newDeleteAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		if stubAddressStore.deletedAddressId != deleteAddressRequest.Id {
			t.Errorf("deleted address with id %d want %d",
				stubAddressStore.deletedAddressId, deleteAddressRequest.Id)
		}
	})
}

func newDeleteAddressRequest(customerJWT string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(http.MethodDelete, "/customer/address", body)
	request.Header.Add("Token", customerJWT)

	return request
}

func TestSaveCustomerAddress(t *testing.T) {
	stubAddressStore := &StubAddressStore{[]Address{}, Address{}, 0}
	stubCustomerStore := &StubCustomerStore{[]Customer{peterCustomer, aliceCustomer}, nil, nil}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	expiresAt := time.Now().Add(time.Second)
	server := CustomerAddressServer{stubAddressStore, stubCustomerStore, secretKey}

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		invalidJWT := "thisIsAnInvalidJWT"
		request := newAddAddressRequest(invalidJWT, nil)
		request.Header.Add("Token", invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Bad Request on inavlid request", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, peterCustomer.Id)
		request := newAddAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns Not Found on missing user", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(aliceAddress)
		missingJWT, _ := GenerateJWT(secretKey, expiresAt, 10)
		request := newAddAddressRequest(missingJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("saves Peter's new address", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(peterAddress1)
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, peterCustomer.Id)
		request := newAddAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		// Copy the address ID from the dummy data to the stored data. This is
		// done because the tests use stubs for store implementations and
		// dummy user data for test cases, so there is going to be a mismatch
		// of the assigned IDs between the two.
		stubAddressStore.savedAddress.Id = peterAddress1.Id
		if !reflect.DeepEqual(stubAddressStore.savedAddress, peterAddress1) {
			t.Errorf("got %v want %v", stubAddressStore.savedAddress, peterAddress1)
		}
	})

	t.Run("saves Alice's new address", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(aliceAddress)
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, aliceCustomer.Id)
		request := newAddAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		// Copy the address ID from the dummy data to the stored data. This is
		// done because the tests use stubs for store implementations and
		// dummy user data for test cases, so there is going to be a mismatch
		// of the assigned IDs between the two.
		stubAddressStore.savedAddress.Id = aliceAddress.Id
		if !reflect.DeepEqual(stubAddressStore.savedAddress, aliceAddress) {
			t.Errorf("got %v want %v", stubAddressStore.savedAddress, aliceAddress)
		}
	})
}

func newAddAddressRequest(customerJWT string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/customer/address", body)
	request.Header.Add("Token", customerJWT)

	return request
}

func TestGetCustomerAddress(t *testing.T) {
	stubAddressStore := &StubAddressStore{
		[]Address{peterAddress1, peterAddress2, aliceAddress}, Address{}, 0,
	}
	stubCustomerStore := &StubCustomerStore{[]Customer{peterCustomer, aliceCustomer}, nil, nil}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	expiresAt := time.Now().Add(time.Second)
	server := CustomerAddressServer{stubAddressStore, stubCustomerStore, secretKey}

	t.Run("returns Peter's addresses", func(t *testing.T) {
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, peterCustomer.Id)
		request := newGetCustomerAddressRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		want := []GetAddressResponse{
			addressToGetAddressResponse(peterAddress1),
			addressToGetAddressResponse(peterAddress2),
		}
		var got []GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)
		assertAddresses(t, got, want)
	})

	t.Run("returns Alice's addresses", func(t *testing.T) {
		aliceJWT, _ := GenerateJWT(secretKey, expiresAt, aliceCustomer.Id)
		request := newGetCustomerAddressRequest(aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		want := []GetAddressResponse{
			addressToGetAddressResponse(aliceAddress),
		}
		var got []GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)

		assertAddresses(t, got, want)
	})

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		invalidJWT := "thisIsAnInvalidJWT"
		request := newGetCustomerAddressRequest(invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Not Found on missing user", func(t *testing.T) {
		aliceJWT, _ := GenerateJWT(secretKey, expiresAt, 10)
		request := newGetCustomerAddressRequest(aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		assertErrorResponse(t, errorResponse, ErrMissingCustomer)
	})
}

func newGetCustomerAddressRequest(customerJWT string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/customer/address/", nil)
	request.Header.Add("Token", customerJWT)

	return request
}

func assertAddresses(t testing.TB, got, want []GetAddressResponse) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, aliceAddress)
	}
}
