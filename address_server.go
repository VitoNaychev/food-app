package bt_customer_svc

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type CustomerAddressStore interface {
	GetAddressesByCustomerId(customerId int) ([]Address, error)
	StoreAddress(address Address)
}

type CustomerAddressServer struct {
	addressStore  CustomerAddressStore
	customerStore CustomerStore
	secretKey     []byte
}

type GetAddressResponse struct {
	Lat          float64 `valid:"latitude,required"`
	Lon          float64 `valid:"longitude,required"`
	AddressLine1 string  `valid:"stringlength(2|40),required"`
	AddressLine2 string  `valid:"stringlength(40)"`
	City         string  `valid:"stringlength(2|20),required"`
	Country      string  `valid:"stringlength(2|20),required"`
}

type AddAddressRequest struct {
	Lat          float64 `valid:"latitude,required"`
	Lon          float64 `valid:"longitude,required"`
	AddressLine1 string  `valid:"stringlength(2|40),required"`
	AddressLine2 string  `valid:"stringlength(40)"`
	City         string  `valid:"stringlength(2|20),required"`
	Country      string  `valid:"stringlength(2|20),required"`
}

type Address struct {
	CustomerId   int
	Lat          float64
	Lon          float64
	AddressLine1 string
	AddressLine2 string
	City         string
	Country      string
}

func (c *CustomerAddressServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AuthenticationMiddleware(c.StoreAddressHandler, c.secretKey)(w, r)
	case http.MethodGet:
		AuthenticationMiddleware(c.GetAddressHandler, c.secretKey)(w, r)
	}
}

func (c *CustomerAddressServer) StoreAddressHandler(w http.ResponseWriter, r *http.Request) {
	var addAddressRequest AddAddressRequest
	err := ValidateBody(r.Body, &addAddressRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		return
	}

	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err = c.customerStore.GetCustomerById(customerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrMissingCustomer.Error()})
		return
	}

	address := addAddressRequestToAddress(addAddressRequest, customerId)

	c.addressStore.StoreAddress(address)
}

func addAddressRequestToAddress(addAddressRequest AddAddressRequest, customerId int) Address {
	address := Address{
		CustomerId:   customerId,
		Lat:          addAddressRequest.Lat,
		Lon:          addAddressRequest.Lon,
		AddressLine1: addAddressRequest.AddressLine1,
		AddressLine2: addAddressRequest.AddressLine2,
		City:         addAddressRequest.City,
		Country:      addAddressRequest.Country,
	}

	return address
}

func (c *CustomerAddressServer) GetAddressHandler(w http.ResponseWriter, r *http.Request) {
	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err := c.customerStore.GetCustomerById(customerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrMissingCustomer.Error()})
		return
	}

	addresses, _ := c.addressStore.GetAddressesByCustomerId(customerId)

	getAddressResponse := []GetAddressResponse{}
	for _, address := range addresses {
		getAddressResponse = append(getAddressResponse, addressToGetAddressResponse(address))
	}

	json.NewEncoder(w).Encode(getAddressResponse)
}

func addressToGetAddressResponse(address Address) GetAddressResponse {
	getAddressResponse := GetAddressResponse{
		Lat:          address.Lat,
		Lon:          address.Lon,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		Country:      address.Country,
	}

	return getAddressResponse
}
