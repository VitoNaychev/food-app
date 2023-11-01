package customer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/VitoNaychev/bt-customer-svc/handlers/auth"
	cs "github.com/VitoNaychev/bt-customer-svc/models/customer_store"
)

func NewUpdateCustomerRequest(customer cs.Customer, SecretKey []byte, ExpiresAt time.Duration) *http.Request {
	body := bytes.NewBuffer([]byte{})
	updateCustomerRequest := UpdateCustomerRequest{
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Password:    customer.Password,
	}
	json.NewEncoder(body).Encode(updateCustomerRequest)

	customerJWT, _ := auth.GenerateJWT(SecretKey, ExpiresAt, customer.Id)

	request, _ := http.NewRequest(http.MethodPut, "/customer/", body)
	request.Header.Add("Token", customerJWT)
	return request
}

func NewLoginRequest(customer cs.Customer) *http.Request {
	loginCustomerRequest := LoginCustomerRequest{
		Email:    customer.Email,
		Password: customer.Password,
	}
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(loginCustomerRequest)

	request, _ := http.NewRequest(http.MethodPost, "/customer/login/", body)
	return request
}

func NewCreateCustomerRequest(customer cs.Customer) *http.Request {
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

func NewGetCustomerRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/customer/", nil)
	request.Header.Add("Token", jwt)

	return request
}
