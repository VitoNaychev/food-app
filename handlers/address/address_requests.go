package address

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/VitoNaychev/bt-customer-svc/models"
)

func NewUpdateAddressRequest(address models.Address, customerJWT string) *http.Request {
	updateAddressRequest := AddressToUpdateAddressRequest(address)

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(updateAddressRequest)

	request, _ := http.NewRequest(http.MethodPut, "/customer/address/", body)
	request.Header.Add("Token", customerJWT)

	return request
}

func NewDeleteAddressRequest(customerJWT string, id int) *http.Request {
	deleteAddressRequest := DeleteAddressRequest{Id: id}
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(deleteAddressRequest)

	request, _ := http.NewRequest(http.MethodDelete, "/customer/address", body)
	request.Header.Add("Token", customerJWT)

	return request
}

func NewAddAddressRequest(customerJWT string, address models.Address) *http.Request {
	addAddressRequest := AddressToAddAddressRequest(address)
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(addAddressRequest)

	request, _ := http.NewRequest(http.MethodPost, "/customer/address", body)
	request.Header.Add("Token", customerJWT)

	return request
}

func NewGetAddressRequest(customerJWT string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/customer/address/", nil)
	request.Header.Add("Token", customerJWT)

	return request
}
