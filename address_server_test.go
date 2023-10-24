package bt_customer_svc

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

type StubAddressStore struct {
	addresses []GetAddressResponse
}

func (s *StubAddressStore) GetAddressesByCustomerId(customerId int) ([]GetAddressResponse, error) {
	if customerId == peterCustomer.Id {
		return []GetAddressResponse{peterAddress1, peterAddress2}, nil
	}

	if customerId == aliceCustomer.Id {
		return []GetAddressResponse{aliceAddress}, nil
	}

	return []GetAddressResponse{}, nil
}

func TestGetCustomerAddress(t *testing.T) {
	stubAddressStore := &StubAddressStore{[]GetAddressResponse{peterAddress1, peterAddress2, aliceAddress}}
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

		want := []GetAddressResponse{peterAddress1, peterAddress2}
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

		want := []GetAddressResponse{aliceAddress}
		var got []GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)
		assertAddresses(t, got, want)
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
