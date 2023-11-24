package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/VitoNaychev/food-app/customer-svc/models"
)

func NewUpdateAddressRequest(customerJWT string, address models.Address) *http.Request {
	updateAddressRequest := AddressToUpdateAddressRequest(address)

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(updateAddressRequest)

	request, _ := http.NewRequest(http.MethodPut, "/customer/address/", body)
	request.Header.Add("Token", customerJWT)

	return request
}

func NewDeleteAddressRequest(customerJWT string, deleteAddressRequest DeleteAddressRequest) *http.Request {
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(deleteAddressRequest)

	request, _ := http.NewRequest(http.MethodDelete, "/customer/address/", body)
	request.Header.Add("Token", customerJWT)

	return request
}

func NewCreateAddressRequest(customerJWT string, address models.Address) *http.Request {
	createAddressRequest := AddressToCreateAddressRequest(address)
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(createAddressRequest)

	request, _ := http.NewRequest(http.MethodPost, "/customer/address/", body)
	request.Header.Add("Token", customerJWT)

	return request
}

func NewGetAddressRequest(customerJWT string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/customer/address/", nil)
	request.Header.Add("Token", customerJWT)

	return request
}
