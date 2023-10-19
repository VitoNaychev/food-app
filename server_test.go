package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

type StubCustomerStore struct {
	customers []Customer
}

func (s *StubCustomerStore) GetCustomerById(id int) (*Customer, error) {
	if len(s.customers) < id {
		return nil, fmt.Errorf("no customer with id %d", id)
	}

	customer := s.customers[id]

	return &customer, nil
}

func (s *StubCustomerStore) GetCustomerByEmail(email string) (*Customer, error) {
	for _, customer := range s.customers {
		if customer.Email == email {
			return &customer, nil
		}
	}

	return nil, fmt.Errorf("no customer with email %v", email)
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

var peterCustomer = Customer{
	Id:          0,
	FirstName:   "Peter",
	LastName:    "Smith",
	PhoneNumber: "+359 88 576 5981",
	Email:       "petesmith@gmail.com",
	Password:    "firefirefire",
}

var aliceCustomer = Customer{
	Id:          1,
	FirstName:   "Alice",
	LastName:    "Johnson",
	PhoneNumber: "+359 88 444 2222",
	Email:       "alicejohn@gmail.com",
	Password:    "helloJohn123",
}

type MalformedRequest struct {
	s string
}

func TestLoginUser(t *testing.T) {
	store := &StubCustomerStore{[]Customer{peterCustomer, aliceCustomer}}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	server := NewCustomerServer(secretKey, time.Now().Add(time.Second), store)

	t.Run("returns JWT on Peter's credentials", func(t *testing.T) {
		request := newLoginRequest(peterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertJWT(t, response.Header(), secretKey, peterCustomer.Id)
	})

	t.Run("returns JWT on Alice's credentials", func(t *testing.T) {
		request := newLoginRequest(aliceCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertJWT(t, response.Header(), secretKey, aliceCustomer.Id)
	})

	t.Run("returns Unauthorized on invalid credentials", func(t *testing.T) {
		incorrectCustomer := peterCustomer
		incorrectCustomer.Password = "passsword123"
		request := newLoginRequest(incorrectCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		assertStatus(t, response.Code, http.StatusUnauthorized)
		assertErrorResponse(t, errorResponse, ErrInvalidCredentials)
	})

	t.Run("returns Unauthorized on missing user", func(t *testing.T) {
		missingCustomer := peterCustomer
		missingCustomer.Email = "notanemail@gmail.com"
		request := newLoginRequest(missingCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		assertStatus(t, response.Code, http.StatusUnauthorized)
		assertErrorResponse(t, errorResponse, ErrMissingCustomer)
	})

	t.Run("returns Bad Request on invalid request field", func(t *testing.T) {
		malformedRequest := LoginCustomerRequest{
			Email:    "thisisnotan.email",
			Password: "password123",
		}
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(malformedRequest)

		request, _ := http.NewRequest(http.MethodPost, "/customer/login/", body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertErrorResponse(t, errorResponse, ErrMalformedRequest)
	})

	t.Run("returns Bad Request on malformed request", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(MalformedRequest{"malformed request"})

		request, _ := http.NewRequest(http.MethodPost, "/customer/login/", body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertErrorResponse(t, errorResponse, ErrMalformedRequest)
	})
}

func newLoginRequest(customer Customer) *http.Request {
	loginCustomerRequest := LoginCustomerRequest{
		Email:    customer.Email,
		Password: customer.Password,
	}
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(loginCustomerRequest)

	request, _ := http.NewRequest(http.MethodPost, "/customer/login/", body)
	return request
}

func TestCreateUser(t *testing.T) {
	store := &StubCustomerStore{}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	server := NewCustomerServer(secretKey, time.Now().Add(time.Second), store)

	t.Run("stores customer on POST", func(t *testing.T) {
		store.Empty()

		request := newCreateCustomerRequest(peterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusAccepted

		assertStatus(t, got, want)

		if len(store.customers) != 1 {
			t.Errorf("got %d calls to StoreCustomer want %d", len(store.customers), 1)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("returns JWT on POST", func(t *testing.T) {
		store.Empty()

		request := newCreateCustomerRequest(peterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusAccepted

		assertStatus(t, got, want)
		assertJWT(t, response.Header(), secretKey, peterCustomer.Id)
	})

	t.Run("returns Bad Request on invalid request field (phone number)", func(t *testing.T) {
		store.Empty()

		malformedCustomer := Customer{
			Id:          0,
			FirstName:   "Peter",
			LastName:    "Smith",
			PhoneNumber: "+359 aa bbbb 834",
			Email:       "petesmith@gmail.com",
			Password:    "$#andsfkasnflkkadf",
		}

		request := newCreateCustomerRequest(malformedCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertErrorResponse(t, errorResponse, ErrMalformedRequest)
	})

	t.Run("returns Bad Request on malformed request", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(MalformedRequest{"malformed request"})

		request, _ := http.NewRequest(http.MethodPost, "/customer/", body)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertErrorResponse(t, errorResponse, ErrMalformedRequest)
	})

	t.Run("return Bad Request on user with same email", func(t *testing.T) {
		store.Empty()

		request := newCreateCustomerRequest(peterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		request = newCreateCustomerRequest(peterCustomer)
		// reinit ResponseRecorder as it allows a
		// one-time only write of the Status Code
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertErrorResponse(t, errorResponse, ErrExistingUser)
	})
}

func assertErrorResponse(t testing.TB, errorResponse ErrorResponse, expetedError error) {
	t.Helper()

	if errorResponse.Message != expetedError.Error() {
		t.Errorf("got error %q want %q", errorResponse.Message, expetedError.Error())
	}
}

func assertJWT(t testing.TB, header http.Header, secretKey []byte, wantId int) {
	t.Helper()

	if header["Token"] == nil {
		t.Fatalf("missing JWT in header")
	}

	token, err := verifyJWT(header["Token"][0], secretKey)
	if err != nil {
		t.Fatalf("error verifying JWT: %v", err)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		t.Fatalf("did not get subject in JWT, expected %v", wantId)
	}

	gotId, err := strconv.Atoi(subject)
	if err != nil {
		t.Fatalf("did not get customer ID for subject, got %v", subject)
	}

	if gotId != wantId {
		t.Errorf("got customer id %v want %v", subject, wantId)
	}
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

func TestGetUser(t *testing.T) {
	store := &StubCustomerStore{[]Customer{peterCustomer, aliceCustomer}}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	server := NewCustomerServer(secretKey, time.Now(), store)

	t.Run("returns Peter's customer information", func(t *testing.T) {
		peterJWT, _ := generateJWT(secretKey, time.Now().Add(time.Second), peterCustomer.Id)
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
		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		assertStatus(t, response.Code, http.StatusUnauthorized)
		assertErrorResponse(t, errorResponse, ErrMissingToken)
	})

	t.Run("returns Not Found on missing customer", func(t *testing.T) {
		noCustomerJWT, _ := generateJWT(secretKey, time.Now().Add(time.Second), 3)
		request := newGetCustomerRequest(noCustomerJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		assertStatus(t, response.Code, http.StatusNotFound)
		assertErrorResponse(t, errorResponse, ErrMissingCustomer)
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
