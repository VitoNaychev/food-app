package address

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	as "github.com/VitoNaychev/bt-customer-svc/models/address_store"
)

func NewUpdateAddressRequest(address as.Address, customerJWT string) *http.Request {
	updateAddressRequest := UpdateAddressRequest{
		Id:           address.Id,
		Lat:          address.Lat,
		Lon:          address.Lon,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		Country:      address.Country,
	}

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(updateAddressRequest)

	request, _ := http.NewRequest(http.MethodPut, "/customer/address/", body)
	request.Header.Add("Token", customerJWT)

	return request
}

func NewDeleteAddressRequest(customerJWT string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(http.MethodDelete, "/customer/address", body)
	request.Header.Add("Token", customerJWT)

	return request
}

func NewAddAddressRequest(customerJWT string, body io.Reader) *http.Request {
	request, _ := http.NewRequest(http.MethodPost, "/customer/address", body)
	request.Header.Add("Token", customerJWT)

	return request
}

func NewGetAddressRequest(customerJWT string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/customer/address/", nil)
	request.Header.Add("Token", customerJWT)

	return request
}
