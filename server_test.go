package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type StubAddressStore struct {
	addresses []GetAddressResponse
}

func (s *StubAddressStore) GetAddressByCustomerId(customerId int) (*GetAddressResponse, error) {
	if len(s.addresses) < customerId {
		return nil, errors.New("customer doesn't have an address")
	}

	return &s.addresses[customerId], nil
}

type StubCustomerStore struct {
	customers   []Customer
	deleteCalls []int
	updateCalls []Customer
}

func (s *StubCustomerStore) UpdateCustomer(customer Customer) error {
	s.updateCalls = append(s.updateCalls, customer)

	return nil
}

func (s *StubCustomerStore) DeleteCustomer(id int) error {
	for _, customer := range s.customers {
		if customer.Id == id {
			// s.customers = append(s.customers[:id], s.customers[id+1:]...)
			s.deleteCalls = append(s.deleteCalls, id)
			return nil
		}
	}

	return fmt.Errorf("no customer with id %d", id)
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

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

type DummyRequest struct {
	S string `valid:"stringlength(10|20),required"`
	I int    `valid:"required"`
}

type IncorrectDummyRequest struct {
	S int
	I string
}

func TestGetCustomerAddress(t *testing.T) {
	stubAddressStore := StubAddressStore{[]GetAddressResponse{peterAddress, aliceAddress}}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	expiresAt := time.Now().Add(time.Second)
	server := NewCustomerServer(secretKey, expiresAt, nil, &stubAddressStore)

	t.Run("returns Peter's addresses", func(t *testing.T) {
		peterJWT, _ := generateJWT(secretKey, expiresAt, peterCustomer.Id)
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
		aliceJWT, _ := generateJWT(secretKey, expiresAt, aliceCustomer.Id)
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

func TestDecodeRequest(t *testing.T) {
	t.Run("returns Bad Request on empty body", func(t *testing.T) {
		body := bytes.NewBuffer([]byte{})

		var dummyRequest DummyRequest
		err := decodeRequest(body, &dummyRequest)

		assertError(t, err, ErrEmptyBody)
	})

	t.Run("returns Bad Request on empty JSON", func(t *testing.T) {
		body := bytes.NewBuffer([]byte(`{}`))

		var dummyRequest DummyRequest
		err := decodeRequest(body, &dummyRequest)

		assertError(t, err, ErrEmptyJSON)
	})

	t.Run("returns Bad Request on incorrect request type", func(t *testing.T) {
		incorrectDummyRequest := IncorrectDummyRequest{
			S: 10,
			I: "Hello, World!",
		}

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(incorrectDummyRequest)

		var dummyRequest DummyRequest
		err := decodeRequest(body, &dummyRequest)

		assertError(t, err, ErrIncorrectRequestType)
	})

	t.Run("returns Bad Request on invalid fields", func(t *testing.T) {
		invalidDummyRequest := DummyRequest{
			S: "Hello,",
			I: 10,
		}

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(invalidDummyRequest)

		var dummyRequest DummyRequest
		err := decodeRequest(body, &dummyRequest)

		assertError(t, err, ErrInvalidRequestField)
	})

	t.Run("returns Accepted on valid request", func(t *testing.T) {
		wantDummyRequest := DummyRequest{
			S: "Hello, World!",
			I: 10,
		}

		body := bytes.NewBuffer([]byte{})
		json.NewEncoder(body).Encode(wantDummyRequest)

		var gotDummyRequest DummyRequest
		err := decodeRequest(body, &gotDummyRequest)

		if err != nil {
			t.Errorf("did not expect error, got %v", err)
		}

		if !reflect.DeepEqual(gotDummyRequest, wantDummyRequest) {
			t.Errorf("got %v want %v", gotDummyRequest, wantDummyRequest)
		}
	})
}

func assertError(t testing.TB, got, want error) {
	t.Helper()

	if got != want {
		t.Errorf("got error %v want %v", got, want)
	}
}

func TestAuthenticationMiddleware(t *testing.T) {
	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	dummyHandler := authenticationMiddleware(DummyHandler, secretKey)

	t.Run("returns Unauthorized on missing JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		assertErrorResponse(t, errorResponse, ErrMissingToken)
	})

	t.Run("returns Unauthorized on invalid JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Token", "thisIsAnInvalidJWT")

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)

		var errorResponse ErrorResponse
		decoder := json.NewDecoder(response.Body)
		decoder.DisallowUnknownFields()
		decoder.Decode(&errorResponse)
		if errorResponse.Message == "" {
			t.Errorf("expected error message but did not get one")
		}
	})

	t.Run("returns Unauthorized on missing Subject in Token", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second)),
		})

		tokenString, _ := token.SignedString(secretKey)
		request.Header.Add("Token", tokenString)

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		assertErrorResponse(t, errorResponse, ErrMissingSubject)
	})

	t.Run("returns Unauthorized on noninteger Subject in Token", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			Subject:   "notAnIntegerSubject",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second)),
		})

		tokenString, _ := token.SignedString(secretKey)
		request.Header.Add("Token", tokenString)

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)
		assertErrorResponse(t, errorResponse, ErrNonIntegerSubject)
	})

	t.Run("returns Token's Subject on valid JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/", nil)
		response := httptest.NewRecorder()

		want := 10
		dummyJWT, _ := generateJWT(secretKey, time.Now().Add(time.Second), want)
		request.Header.Add("Token", dummyJWT)

		dummyHandler(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		subject := request.Header["Subject"]
		if subject == nil {
			t.Fatalf("did not get Subject in request header")
		}

		got, err := strconv.Atoi(subject[0])
		if err != nil {
			t.Fatalf("expected integer Subject, got %q", subject[0])
		}

		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	store := &StubCustomerStore{[]Customer{peterCustomer, aliceCustomer}, []int{}, []Customer{}}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	expiresAt := time.Now().Add(time.Second)
	server := NewCustomerServer(secretKey, expiresAt, store, nil)

	t.Run("updates customer information on valid JWT", func(t *testing.T) {
		customer := peterCustomer
		customer.FirstName = "John"
		customer.PhoneNumber = "+359 88 1234 213"

		request := newUpdateCustomerRequest(customer, secretKey, expiresAt)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

		if len(store.updateCalls) != 1 {
			t.Fatalf("got %d calls to UpdateCustomer expected %d", len(store.updateCalls), 1)
		}

		if !reflect.DeepEqual(store.updateCalls[0], customer) {
			t.Errorf("did not update correct customer got %v want %v", store.updateCalls[0], customer)
		}
	})
}

func newUpdateCustomerRequest(customer Customer, secretKey []byte, expiresAt time.Time) *http.Request {
	body := bytes.NewBuffer([]byte{})
	updateCustomerRequest := UpdateCustomerRequest{
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Password:    customer.Password,
	}
	json.NewEncoder(body).Encode(updateCustomerRequest)

	customerJWT, _ := generateJWT(secretKey, expiresAt, customer.Id)

	request, _ := http.NewRequest(http.MethodPut, "/customer/", body)
	request.Header.Add("Token", customerJWT)
	return request
}

func TestDeleteUser(t *testing.T) {
	store := &StubCustomerStore{[]Customer{peterCustomer, aliceCustomer}, []int{}, []Customer{}}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	expiresAt := time.Now().Add(time.Second)
	server := NewCustomerServer(secretKey, expiresAt, store, nil)

	t.Run("deletes customer on valid JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/customer/", nil)
		response := httptest.NewRecorder()

		peterJWT, _ := generateJWT(secretKey, expiresAt, peterCustomer.Id)
		request.Header.Add("Token", peterJWT)

		server.ServeHTTP(response, request)

		if len(store.deleteCalls) != 1 {
			t.Fatalf("got %d calls to DeleteCustomer expected %d", len(store.deleteCalls), 1)
		}

		if store.deleteCalls[0] != peterCustomer.Id {
			t.Errorf("did not delete correct customer got %d want %d", store.deleteCalls[0], peterCustomer.Id)
		}
	})

	t.Run("returns Not Found on missing customer", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodDelete, "/customer/", nil)
		response := httptest.NewRecorder()

		missingCustomerJWT, _ := generateJWT(secretKey, expiresAt, 10)
		request.Header.Add("Token", missingCustomerJWT)

		server.ServeHTTP(response, request)

		var errorResponse ErrorResponse
		json.NewDecoder(response.Body).Decode(&errorResponse)

		assertStatus(t, response.Code, http.StatusNotFound)
		assertErrorResponse(t, errorResponse, ErrMissingCustomer)
	})
}

func TestLoginUser(t *testing.T) {
	store := &StubCustomerStore{[]Customer{peterCustomer, aliceCustomer}, []int{}, []Customer{}}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	expiresAt := time.Now().Add(time.Second)
	server := NewCustomerServer(secretKey, expiresAt, store, nil)

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
	expiresAt := time.Now().Add(time.Second)
	server := NewCustomerServer(secretKey, expiresAt, store, nil)

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
	store := &StubCustomerStore{[]Customer{peterCustomer, aliceCustomer}, []int{}, []Customer{}}

	godotenv.Load("test.env")
	secretKey := []byte(os.Getenv("SECRET"))
	expiresAt := time.Now().Add(time.Second)
	server := NewCustomerServer(secretKey, expiresAt, store, nil)

	t.Run("returns Peter's customer information", func(t *testing.T) {
		peterJWT, _ := generateJWT(secretKey, expiresAt, peterCustomer.Id)
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
		aliceJWT, _ := generateJWT(secretKey, expiresAt, 1)
		request := newGetCustomerRequest(aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var got GetCustomerResponse
		json.NewDecoder(response.Body).Decode(&got)

		customer := store.customers[1]

		assertStatus(t, response.Code, http.StatusOK)
		assertGetCustomerResponse(t, got, customer)
	})

	t.Run("returns Not Found on missing customer", func(t *testing.T) {
		noCustomerJWT, _ := generateJWT(secretKey, expiresAt, 3)
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
