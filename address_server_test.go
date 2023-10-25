package bt_customer_svc

import (
	"bytes"
	"encoding/json"
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
	addresses    []Address
	savedAddress Address
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

func TestSaveCustomerAddress(t *testing.T) {
	stubAddressStore := &StubAddressStore{[]Address{}, Address{}}
	stubCustomerStore := &StubCustomerStore{[]Customer{peterCustomer, aliceCustomer}, nil, nil}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	expiresAt := time.Now().Add(time.Second)
	server := CustomerAddressServer{stubAddressStore, stubCustomerStore, secretKey}

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		invalidJWT := "thisIsAnInvalidJWT"
		request, _ := http.NewRequest(http.MethodPost, "/customer/address", nil)
		request.Header.Add("Token", invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Bad Request on inavlid request", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		request, _ := http.NewRequest(http.MethodPost, "/customer/address", body)
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, peterCustomer.Id)
		request.Header.Add("Token", peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})

	t.Run("returns Not Found on missing user", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(peterAddress2)
		request, _ := http.NewRequest(http.MethodPost, "/customer/address", body)
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, 10)
		request.Header.Add("Token", peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("saves Peter's new address", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(peterAddress1)
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, peterCustomer.Id)
		request := newStoreAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		if !reflect.DeepEqual(stubAddressStore.savedAddress, peterAddress1) {
			t.Errorf("got %v want %v", stubAddressStore.savedAddress, peterAddress1)
		}
	})

	t.Run("saves Alice's new address", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(aliceAddress)
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, aliceCustomer.Id)
		request := newStoreAddressRequest(peterJWT, body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		if !reflect.DeepEqual(stubAddressStore.savedAddress, aliceAddress) {
			t.Errorf("got %v want %v", stubAddressStore.savedAddress, aliceAddress)
		}
	})
}

func newStoreAddressRequest(customerJWT string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/customer/address", body)
	request.Header.Add("Token", customerJWT)

	return request
}

func TestGetCustomerAddress(t *testing.T) {
	stubAddressStore := &StubAddressStore{[]Address{peterAddress1, peterAddress2, aliceAddress}, Address{}}
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
