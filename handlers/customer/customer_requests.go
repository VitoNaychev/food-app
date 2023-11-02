package customer

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/VitoNaychev/bt-customer-svc/handlers/auth"
	"github.com/VitoNaychev/bt-customer-svc/models"
)

func NewUpdateCustomerRequest(customer models.Customer, SecretKey []byte, ExpiresAt time.Duration) *http.Request {
	body := bytes.NewBuffer([]byte{})
	updateCustomerRequest := CustomerToUpdateCustomerRequest(customer)
	json.NewEncoder(body).Encode(updateCustomerRequest)

	customerJWT, _ := auth.GenerateJWT(SecretKey, ExpiresAt, customer.Id)

	request, _ := http.NewRequest(http.MethodPut, "/customer/", body)
	request.Header.Add("Token", customerJWT)
	return request
}

func NewLoginRequest(customer models.Customer) *http.Request {
	loginCustomerRequest := CustomerToLoginCustomerRequest(customer)
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(loginCustomerRequest)

	request, _ := http.NewRequest(http.MethodPost, "/customer/login/", body)
	return request
}

func NewCreateCustomerRequest(customer models.Customer) *http.Request {
	createCustomerRequest := CustomerToCreateCustomerRequest(customer)
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
