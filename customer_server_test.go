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
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

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

func DummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
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

	customerJWT, _ := GenerateJWT(secretKey, expiresAt, customer.Id)

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

		peterJWT, _ := GenerateJWT(secretKey, expiresAt, peterCustomer.Id)
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

		missingCustomerJWT, _ := GenerateJWT(secretKey, expiresAt, 10)
		request.Header.Add("Token", missingCustomerJWT)

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
		assertErrorResponse(t, response.Body, ErrMissingCustomer)
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

		assertStatus(t, response.Code, http.StatusUnauthorized)
		assertErrorResponse(t, response.Body, ErrInvalidCredentials)
	})

	t.Run("returns Unauthorized on missing user", func(t *testing.T) {
		missingCustomer := peterCustomer
		missingCustomer.Email = "notanemail@gmail.com"
		request := newLoginRequest(missingCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusUnauthorized)
		assertErrorResponse(t, response.Body, ErrMissingCustomer)
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

		assertStatus(t, response.Code, http.StatusBadRequest)
		assertErrorResponse(t, response.Body, ErrExistingUser)
	})
}

func assertErrorResponse(t testing.TB, body io.Reader, expetedError error) {
	t.Helper()

	var errorResponse ErrorResponse
	json.NewDecoder(body).Decode(&errorResponse)

	if errorResponse.Message != expetedError.Error() {
		t.Errorf("got error %q want %q", errorResponse.Message, expetedError.Error())
	}
}

func assertJWT(t testing.TB, header http.Header, secretKey []byte, wantId int) {
	t.Helper()

	if header["Token"] == nil {
		t.Fatalf("missing JWT in header")
	}

	token, err := VerifyJWT(header["Token"][0], secretKey)
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
		peterJWT, _ := GenerateJWT(secretKey, expiresAt, peterCustomer.Id)
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
		aliceJWT, _ := GenerateJWT(secretKey, expiresAt, 1)
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
		noCustomerJWT, _ := GenerateJWT(secretKey, expiresAt, 3)
		request := newGetCustomerRequest(noCustomerJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
		assertErrorResponse(t, response.Body, ErrMissingCustomer)
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
		t.Fatalf("got status %d, want %d", got, want)
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
