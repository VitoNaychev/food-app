package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/customer-svc/handlers"
	"github.com/VitoNaychev/food-app/customer-svc/models"
	td "github.com/VitoNaychev/food-app/customer-svc/testdata"
	"github.com/VitoNaychev/food-app/customer-svc/testutil"
	"github.com/VitoNaychev/food-app/msgtypes"
	"github.com/VitoNaychev/food-app/validation"
	"github.com/golang-jwt/jwt/v5"
)

func TestCustomerEndpointAuthentication(t *testing.T) {
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := handlers.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	invalidJWT := "thisIsAnInvalidJWT"
	cases := map[string]*http.Request{
		"get customer authentication":    handlers.NewGetCustomerRequest(invalidJWT),
		"update customer authentication": handlers.NewDeleteCustomerRequest(invalidJWT),
		"delete customer authentication": handlers.NewDeleteCustomerRequest(invalidJWT),
	}

	for name, request := range cases {
		t.Run(name, func(t *testing.T) {
			request.Header.Add("Token", invalidJWT)

			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		})
	}
}

func TestAuthHandler(t *testing.T) {
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := handlers.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("returns OK status and customer ID on valid JWT", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := handlers.NewAuthRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := msgtypes.AuthResponse{
			Status: msgtypes.OK,
			ID:     td.PeterCustomer.Id,
		}
		var got msgtypes.AuthResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns INVALID status on invalid JWT", func(t *testing.T) {
		invalidJWT := "invalidJWT"

		request := handlers.NewAuthRequest(invalidJWT)
		request.Header.Add("Token", invalidJWT)

		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := msgtypes.AuthResponse{
			Status: msgtypes.INVALID,
			ID:     0,
		}
		var got msgtypes.AuthResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns MISSING_TOKEN status on missing JWT", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/customer/auth/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := msgtypes.AuthResponse{
			Status: msgtypes.MISSING_TOKEN,
			ID:     0,
		}
		var got msgtypes.AuthResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns INVALID on noninteger subject", func(t *testing.T) {
		invalidJWT, _ := GenerateJWTWithStringSubject(testEnv.SecretKey, testEnv.ExpiresAt, "peter")
		request := handlers.NewAuthRequest(invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := msgtypes.AuthResponse{
			Status: msgtypes.INVALID,
			ID:     0,
		}
		var got msgtypes.AuthResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns INVALID on missing subject", func(t *testing.T) {
		invalidJWT, _ := GenerateJWTWithoutSubject(testEnv.SecretKey, testEnv.ExpiresAt)
		request := handlers.NewAuthRequest(invalidJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := msgtypes.AuthResponse{
			Status: msgtypes.INVALID,
			ID:     0,
		}
		var got msgtypes.AuthResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns NOT_FOUND on customer that doesn't exist", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)

		request := handlers.NewAuthRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		want := msgtypes.AuthResponse{
			Status: msgtypes.NOT_FOUND,
			ID:     0,
		}
		var got msgtypes.AuthResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

}

func TestUpdateUser(t *testing.T) {
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := handlers.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("updates customer information on valid JWT", func(t *testing.T) {
		updateCustomer := td.PeterCustomer
		updateCustomer.FirstName = "John"
		updateCustomer.PhoneNumber = "+359 88 1234 213"

		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := handlers.NewUpdateCustomerRequest(updateCustomer, peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertUpdatedCustomer(t, store, updateCustomer)
	})
}

func TestDeleteUser(t *testing.T) {
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := handlers.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("deletes customer on valid JWT", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)

		request := handlers.NewDeleteCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		testutil.AssertDeletedCustomer(t, store, td.PeterCustomer)
	})

	t.Run("returns Not Found on missing customer", func(t *testing.T) {
		missingCustomerJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)

		request := handlers.NewDeleteCustomerRequest(missingCustomerJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrCustomerNotFound)
	})
}

func TestLoginUser(t *testing.T) {
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := handlers.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("returns JWT on Peter's credentials", func(t *testing.T) {
		request := handlers.NewLoginRequest(td.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)

		var jwtResponse handlers.JWTResponse
		json.NewDecoder(response.Body).Decode(&jwtResponse)

		testutil.AssertJWT(t, jwtResponse.Token, testEnv.SecretKey, td.PeterCustomer.Id)
	})

	t.Run("returns JWT on Alice's credentials", func(t *testing.T) {
		request := handlers.NewLoginRequest(td.AliceCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)

		var jwtResponse handlers.JWTResponse
		json.NewDecoder(response.Body).Decode(&jwtResponse)

		testutil.AssertJWT(t, jwtResponse.Token, testEnv.SecretKey, td.AliceCustomer.Id)
	})

	t.Run("returns Unauthorized on invalid credentials", func(t *testing.T) {
		incorrectCustomer := td.PeterCustomer
		incorrectCustomer.Password = "passsword123"
		request := handlers.NewLoginRequest(incorrectCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidCredentials)
	})

	t.Run("returns Unauthorized on missing user", func(t *testing.T) {
		missingCustomer := td.PeterCustomer
		missingCustomer.Email = "notanemail@gmail.com"
		request := handlers.NewLoginRequest(missingCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrCustomerNotFound)
	})
}

func TestCreateUser(t *testing.T) {
	customerData := []models.Customer{}
	store := testutil.NewStubCustomerStore(customerData)
	server := handlers.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("stores customer on POST", func(t *testing.T) {
		store.Empty()

		request := handlers.NewCreateCustomerRequest(td.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := response.Code
		want := http.StatusAccepted

		testutil.AssertStatus(t, got, want)
		testutil.AssertCreatedCustomer(t, store, td.PeterCustomer)

	})

	t.Run("returns JSON with JWT and new customer on POST", func(t *testing.T) {
		store.Empty()

		request := handlers.NewCreateCustomerRequest(td.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusAccepted)

		wantResponseCustomer := handlers.CustomerToCustomerResponse(td.PeterCustomer)

		var gotResponse handlers.CreateCustomerResponse
		json.NewDecoder(response.Body).Decode(&gotResponse)

		testutil.AssertJWT(t, gotResponse.JWT.Token, testEnv.SecretKey, td.PeterCustomer.Id)
		testutil.AssertEqual(t, gotResponse.Customer, wantResponseCustomer)
	})

	t.Run("return Bad Request on user with same email", func(t *testing.T) {
		store.Empty()

		request := handlers.NewCreateCustomerRequest(td.PeterCustomer)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		request = handlers.NewCreateCustomerRequest(td.PeterCustomer)
		// reinit ResponseRecorder as it allows a
		// one-time only write of the Status Code
		response = httptest.NewRecorder()
		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrExistingCustomer)
	})
}

func TestGetUser(t *testing.T) {
	customerData := []models.Customer{td.PeterCustomer, td.AliceCustomer}
	store := testutil.NewStubCustomerStore(customerData)
	server := handlers.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("returns Peter's customer information", func(t *testing.T) {
		peterJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.PeterCustomer.Id)
		request := handlers.NewGetCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[handlers.GetCustomerResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		want := handlers.CustomerToGetCustomerResponse(td.PeterCustomer)
		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns Alice's customer information", func(t *testing.T) {
		aliceJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.AliceCustomer.Id)
		request := handlers.NewGetCustomerRequest(aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[handlers.GetCustomerResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		want := handlers.CustomerToGetCustomerResponse(td.AliceCustomer)
		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns Not Found on missing customer", func(t *testing.T) {
		noCustomerJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 3)
		request := handlers.NewGetCustomerRequest(noCustomerJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrCustomerNotFound)
	})
}

func GenerateJWTWithStringSubject(secretKey []byte, expiresAt time.Duration, subject string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateJWTWithoutSubject(secretKey []byte, expiresAt time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),
	})

	tokenString, err := token.SignedString(secretKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
