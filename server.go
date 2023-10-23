package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/golang-jwt/jwt/v5"
)

type CustomerStore interface {
	GetCustomerById(id int) (*Customer, error)
	GetCustomerByEmail(email string) (*Customer, error)
	StoreCustomer(customer Customer) int
	DeleteCustomer(id int) error
	UpdateCustomer(customer Customer) error
}

type CustomerServer struct {
	secretKey []byte
	expiresAt time.Time
	store     CustomerStore
	http.Handler
}

func NewCustomerServer(secretKey []byte, expiresAt time.Time, store CustomerStore) *CustomerServer {
	govalidator.TagMap["phonenumber"] = govalidator.Validator(func(str string) bool {
		matched, _ := regexp.Match(`^\+[\d ]+$`, []byte(str))
		return matched
	})

	c := new(CustomerServer)

	c.secretKey = secretKey
	c.expiresAt = expiresAt
	c.store = store

	router := http.NewServeMux()
	router.HandleFunc("/customer/", c.CustomerHandler)
	router.HandleFunc("/customer/login/", c.LoginHandler)

	c.Handler = router

	return c
}

type Customer struct {
	Id          int
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Password    string
}

type GetCustomerResponse struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
}

type CreateCustomerRequest struct {
	FirstName   string `valid:"stringlength(2|20),required"`
	LastName    string `valid:"stringlength(2|20),required"`
	PhoneNumber string `valid:"phonenumber,required"`
	Email       string `valid:"email,required"`
	Password    string `valid:"stringlength(2|72),required"`
}

type LoginCustomerRequest struct {
	Email    string `valid:"email,required"`
	Password string `valid:"stringlength(2|72),required"`
}

type UpdateCustomerRequest struct {
	FirstName   string `valid:"stringlength(2|20),required"`
	LastName    string `valid:"stringlength(2|20),required"`
	PhoneNumber string `valid:"phonenumber,required"`
	Email       string `valid:"email,required"`
	Password    string `valid:"stringlength(2|72),required"`
}

type ErrorResponse struct {
	Message string
}

var (
	ErrExistingUser         = errors.New("user with this email already exists")
	ErrMissingCustomer      = errors.New("customer doesn't exists")
	ErrMissingToken         = errors.New("missing token")
	ErrInvalidCredentials   = errors.New("invalid user credentials")
	ErrMissingSubject       = errors.New("token does not contain subject field")
	ErrNonIntegerSubject    = errors.New("token subject field is not an integer")
	ErrMissingBody          = errors.New("request body is empty")
	ErrEmptyJSON            = errors.New("request JSON is empty")
	ErrIncorrectRequestType = errors.New("request type is incorrect")
	ErrInvalidRequestField  = errors.New("request contains invalid field(s)")
)

func validationMiddleware(endpointHandler func(w http.ResponseWriter, r *http.Request), requestType reflect.Type) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var maxRequestSize int64 = 10000
		body, err := io.ReadAll(io.LimitReader(r.Body, maxRequestSize))
		if string(body) == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: ErrMissingBody.Error()})
			return
		}

		if string(body) == "{}" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: ErrEmptyJSON.Error()})
			return
		}

		requestObject := reflect.New(requestType)

		err = json.Unmarshal(body, requestObject.Interface())
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: ErrIncorrectRequestType.Error()})
			return
		}

		valid, _ := govalidator.ValidateStruct(requestObject.Interface())
		if !valid {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Message: ErrInvalidRequestField.Error()})
			return
		}

		r.Body = io.NopCloser(bytes.NewReader(body))
		endpointHandler(w, r)
	})
}

func authenticationMiddleware(endpointHandler func(w http.ResponseWriter, r *http.Request), secretKey []byte) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Message: ErrMissingToken.Error()})
			return
		}

		token, err := verifyJWT(r.Header["Token"][0], secretKey)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
			return
		}

		id, err := getIDFromToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
			return
		}

		r.Header.Add("Subject", strconv.Itoa(id))

		endpointHandler(w, r)
	})
}

func (c *CustomerServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginCustomerRequest LoginCustomerRequest
	json.NewDecoder(r.Body).Decode(&loginCustomerRequest)

	valid, _ := govalidator.ValidateStruct(loginCustomerRequest)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrInvalidRequestField.Error()})
		return
	}

	customer, err := c.store.GetCustomerByEmail(loginCustomerRequest.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrMissingCustomer.Error()})
		return
	}

	if customer.Password != loginCustomerRequest.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrInvalidCredentials.Error()})
		return
	}

	loginJWT, _ := generateJWT(c.secretKey, c.expiresAt, customer.Id)

	w.WriteHeader(http.StatusAccepted)
	w.Header().Add("Token", loginJWT)
}

func (c *CustomerServer) CustomerHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.storeCustomer(w, r)
	case http.MethodGet:
		authenticationMiddleware(c.getCustomer, c.secretKey)(w, r)
	case http.MethodDelete:
		authenticationMiddleware(c.deleteCustomer, c.secretKey)(w, r)
	case http.MethodPut:
		authenticationMiddleware(c.updateCustomer, c.secretKey)(w, r)
	}
}

func (c *CustomerServer) updateCustomer(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Header["Subject"][0])
	customer, _ := c.store.GetCustomerById(id)

	var updateCustomerRequest UpdateCustomerRequest
	json.NewDecoder(r.Body).Decode(&updateCustomerRequest)
	valid, _ := govalidator.ValidateStruct(updateCustomerRequest)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrInvalidRequestField.Error()})
		return
	}

	customer.FirstName = updateCustomerRequest.FirstName
	customer.LastName = updateCustomerRequest.LastName
	customer.Email = updateCustomerRequest.Email
	customer.PhoneNumber = updateCustomerRequest.PhoneNumber
	customer.Password = updateCustomerRequest.Password

	c.store.UpdateCustomer(*customer)
}

func (c *CustomerServer) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Header["Subject"][0])
	err := c.store.DeleteCustomer(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrMissingCustomer.Error()})
	}
}

func (c *CustomerServer) storeCustomer(w http.ResponseWriter, r *http.Request) {
	var createCustomerRequest CreateCustomerRequest
	json.NewDecoder(r.Body).Decode(&createCustomerRequest)

	valid, _ := govalidator.ValidateStruct(createCustomerRequest)
	if !valid {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrInvalidRequestField.Error()})
		return
	}

	if _, err := c.store.GetCustomerByEmail(createCustomerRequest.Email); err == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrExistingUser.Error()})
		return
	}

	customer := getCustomerFromCreateCustomerRequest(createCustomerRequest)
	customerId := c.store.StoreCustomer(*customer)

	customerJWT, _ := generateJWT(c.secretKey, c.expiresAt, customerId)

	w.WriteHeader(http.StatusAccepted)
	w.Header().Add("Token", customerJWT)
}

func getCustomerFromCreateCustomerRequest(createCustomerRequest CreateCustomerRequest) *Customer {
	customer := &Customer{
		Id:          0,
		FirstName:   createCustomerRequest.FirstName,
		LastName:    createCustomerRequest.LastName,
		PhoneNumber: createCustomerRequest.PhoneNumber,
		Email:       createCustomerRequest.Email,
		Password:    createCustomerRequest.Password,
	}
	return customer
}

func (c *CustomerServer) getCustomer(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Header["Subject"][0])
	customer, err := c.store.GetCustomerById(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrMissingCustomer.Error()})
		return
	}

	getCustomerResponse := newGetCustomerResponse(*customer)
	json.NewEncoder(w).Encode(getCustomerResponse)

}

func newGetCustomerResponse(customer Customer) *GetCustomerResponse {
	getCustomerResponse := GetCustomerResponse{
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		PhoneNumber: customer.PhoneNumber,
		Email:       customer.Email,
	}

	return &getCustomerResponse
}

func getIDFromToken(token *jwt.Token) (int, error) {
	subject, err := token.Claims.GetSubject()
	if err != nil || subject == "" {
		return -1, ErrMissingSubject
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return -1, ErrNonIntegerSubject
	}

	return id, nil
}
