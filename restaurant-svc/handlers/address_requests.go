package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

func NewUpdateAddressRequest(restaurantJWT string, address models.Address) *http.Request {
	updateAddressRequest := AddressToUpdateAddressRequest(address)

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(updateAddressRequest)

	request, _ := http.NewRequest(http.MethodPut, "/restaurant/address/", body)
	request.Header.Add("Token", restaurantJWT)

	return request
}

func NewCreateAddressRequest(restaurantJWT string, address models.Address) *http.Request {
	createAddressRequest := AddressToCreateAddressRequest(address)
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(createAddressRequest)

	request, _ := http.NewRequest(http.MethodPost, "/restaurant/address/", body)
	request.Header.Add("Token", restaurantJWT)

	return request
}

func NewGetAddressRequest(restaurantJWT string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/restaurant/address/", nil)
	request.Header.Add("Token", restaurantJWT)

	return request
}
