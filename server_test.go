package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

type StubCustomerStore struct {
	customers []Customer
}

func (s *StubCustomerStore) GetCustomer(id int) (*Customer, error) {
	if len(s.customers) < id {
		return nil, fmt.Errorf("no customer with id %v", id)
	}

	customer := s.customers[id]

	return &customer, nil
}

func (s *StubCustomerStore) StoreCustomer(customer Customer) int {
	id := len(s.customers)

	customer.Id = id
	s.customers = append(s.customers, customer)

	return id
}

func (s *StubCustomerStore) Empty() {
	s.customers = []Customer{}
}

func TestCreateUser(t *testing.T) {
	store := &StubCustomerStore{}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	server := &CustomerServer{secretKey, time.Now().Add(time.Second), store}

	customer := Customer{
		Id:          0,
		FirstName:   "Peter",
		LastName:    "Smith",
		PhoneNumber: "+359 88 576 5981",
		Email:       "petesmith@gmail.com",
		Password:    "firefirefire",
	}

	t.Run("stores customer on POST", func(t *testing.T) {
		store.Empty()

		request := newCreateCustomerRequest(customer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusAccepted

		assertStatus(t, got, want)
		assertCustomer(t, store.customers[0], customer)
	})

	t.Run("returns JWT on POST", func(t *testing.T) {
		store.Empty()

		request := newCreateCustomerRequest(customer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusAccepted

		assertStatus(t, got, want)

		customerJWT := response.Header()["Token"]
		if customerJWT != nil {
			request, _ := http.NewRequest(http.MethodGet, "/customer/", nil)
			request.Header.Add("Token", customerJWT[0])
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			var getCustomerResponse GetCustomerResponse
			json.NewDecoder(response.Body).Decode(&getCustomerResponse)

			assertGetCustomerResponse(t, getCustomerResponse, customer)
		} else {
			t.Errorf("didn't recieve JWT for the new customer")
		}
	})
}

func newCreateCustomerRequest(customer Customer) *http.Request {
	createCustomerRequest := CreateCustomerRequest{
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		PhoneNumber: customer.PhoneNumber,
		Email:       customer.Email,
		Password:    customer.Password,
	}
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(createCustomerRequest)

	request, _ := http.NewRequest(http.MethodPost, "/customer/", body)
	return request
}

func assertCustomer(t testing.TB, got, want Customer) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetUser(t *testing.T) {
	store := &StubCustomerStore{
		customers: []Customer{
			{
				Id:          0,
				FirstName:   "Peter",
				LastName:    "Smith",
				PhoneNumber: "+359 88 576 5981",
				Email:       "petesmith@gmail.com",
				Password:    "helloWorld123",
			},
			{
				Id:          1,
				FirstName:   "Alice",
				LastName:    "Johnson",
				PhoneNumber: "+359 88 444 2222",
				Email:       "alicejohn@gmail.com",
				Password:    "helloJohn123",
			},
		},
	}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	server := &CustomerServer{secretKey, time.Now(), store}

	t.Run("returns Peter's customer information", func(t *testing.T) {
		peterJWT, _ := generateJWT(secretKey, time.Now().Add(time.Second), 0)
		request := newGetCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got GetCustomerResponse
		json.NewDecoder(response.Body).Decode(&got)

		customer := store.customers[0]

		assertStatus(t, response.Code, http.StatusOK)
		assertGetCustomerResponse(t, got, customer)
	})

	t.Run("returns Alice's customer information", func(t *testing.T) {
		aliceJWT, _ := generateJWT(secretKey, time.Now().Add(time.Second), 1)
		request := newGetCustomerRequest(aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got GetCustomerResponse
		json.NewDecoder(response.Body).Decode(&got)

		customer := store.customers[1]

		assertStatus(t, response.Code, http.StatusOK)
		assertGetCustomerResponse(t, got, customer)
	})

	t.Run("returns Unauthorized on expired JWT", func(t *testing.T) {
		expiredJWT, _ := generateJWT(secretKey, time.Now(), 0)
		request := newGetCustomerRequest(expiredJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)

	})

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		invalidJWT := "thisIsAnInvalidJWT"
		request := newGetCustomerRequest(invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Unauthorized on missing JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/customer/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
	})

	t.Run("returns Not Found on nonexistent customer", func(t *testing.T) {
		noCustomerJWT, _ := generateJWT(secretKey, time.Now().Add(time.Second), 3)
		request := newGetCustomerRequest(noCustomerJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func newGetCustomerRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/customer/", nil)
	request.Header.Add("Token", jwt)

	return request
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Errorf("got status %d, want %d", got, want)
	}
}

func assertGetCustomerResponse(t testing.TB, got GetCustomerResponse, customer Customer) {
	t.Helper()

	want := GetCustomerResponse{
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		PhoneNumber: customer.PhoneNumber,
		Email:       customer.Email,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}
