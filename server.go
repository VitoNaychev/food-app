package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomerStore interface {
	GetCustomer(id int) (*Customer, error)
	StoreCustomer(customer Customer) int
}

type CustomerServer struct {
	secretKey []byte
	expiresAt time.Time
	store     CustomerStore
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
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Password    string
}

type ErrorResponse struct {
	Message string
}

func (c *CustomerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		c.storeCustomer(w, r)
	case http.MethodGet:
		c.getCustomer(w, r)
	}
}

func (c *CustomerServer) storeCustomer(w http.ResponseWriter, r *http.Request) {
	var createCustomerRequest CreateCustomerRequest
	json.NewDecoder(r.Body).Decode(&createCustomerRequest)

	customer := Customer{
		Id:          0,
		FirstName:   createCustomerRequest.FirstName,
		LastName:    createCustomerRequest.LastName,
		PhoneNumber: createCustomerRequest.PhoneNumber,
		Email:       createCustomerRequest.Email,
		Password:    createCustomerRequest.Password,
	}
	customerId := c.store.StoreCustomer(customer)

	customerJWT, _ := generateJWT(c.secretKey, c.expiresAt, customerId)

	w.WriteHeader(http.StatusAccepted)
	w.Header().Add("Token", customerJWT)
}

func (c *CustomerServer) getCustomer(w http.ResponseWriter, r *http.Request) {
	if token, err := verifyJWT(r.Header, c.secretKey); err == nil {
		id := getIDFromToken(token)
		customer, err := c.store.GetCustomer(id)

		if err == nil {
			getCustomerResponse := GetCustomerResponse{
				FirstName:   customer.FirstName,
				LastName:    customer.LastName,
				PhoneNumber: customer.PhoneNumber,
				Email:       customer.Email,
			}
			json.NewEncoder(w).Encode(getCustomerResponse)
		} else {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
	}
}

func getIDFromToken(token *jwt.Token) int {
	subject, _ := token.Claims.GetSubject()

	id, _ := strconv.Atoi(subject)

	return id
}
