package servers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/VitoNaychev/bt-customer-svc/auth"
	"github.com/joho/godotenv"
)

type StubAddressStore struct {
	addresses []GetAddressResponse
}

func (s *StubAddressStore) GetAddressByCustomerId(customerId int) (*GetAddressResponse, error) {
	if len(s.addresses) <= customerId {
		return nil, errors.New("customer doesn't have an address")
	}

	return &s.addresses[customerId], nil
}

var peterAddress = GetAddressResponse{
	Lat:          42.695111,
	Lon:          23.329184,
	AddressLine1: "ulitsa Gerogi S. Rakovski 96",
	AddressLine2: "",
	City:         "Sofia",
	Country:      "Bulgaria",
}

var aliceAddress = GetAddressResponse{
	Lat:          42.6931204,
	Lon:          23.3225465,
	AddressLine1: "ut. Angel Kanchev 1",
	AddressLine2: "",
	City:         "Sofia",
	Country:      "Bulgaria",
}

func TestGetCustomerAddress(t *testing.T) {
	stubAddressStore := StubAddressStore{[]GetAddressResponse{peterAddress, aliceAddress}}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	expiresAt := time.Now().Add(time.Second)
	server := CustomerAddressServer{&stubAddressStore, secretKey}

	t.Run("returns Peter's addresses", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(secretKey, expiresAt, peterCustomer.Id)
		request, _ := http.NewRequest(http.MethodGet, "/customer/address/", nil)
		request.Header.Add("Token", peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		var got GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)
		if !reflect.DeepEqual(peterAddress, got) {
			t.Errorf("got %v want %v", got, peterAddress)
		}
	})

	t.Run("returns Alice's addresses", func(t *testing.T) {
		aliceJWT, _ := auth.GenerateJWT(secretKey, expiresAt, aliceCustomer.Id)
		request, _ := http.NewRequest(http.MethodGet, "/customer/address/", nil)
		request.Header.Add("Token", aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		var got GetAddressResponse
		json.NewDecoder(response.Body).Decode(&got)
		if !reflect.DeepEqual(aliceAddress, got) {
			t.Errorf("got %v want %v", got, aliceAddress)
		}
	})
}
