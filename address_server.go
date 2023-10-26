package bt_customer_svc

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type CustomerAddressStore interface {
	GetAddressesByCustomerId(customerId int) ([]Address, error)
	StoreAddress(address Address)
	DeleteAddressById(id int) error
	GetAddressById(id int) (Address, error)
}

type CustomerAddressServer struct {
	addressStore  CustomerAddressStore
	customerStore CustomerStore
	secretKey     []byte
}

type DeleteAddressRequest struct {
	Id int `validate:"min=0"`
}

type GetAddressResponse struct {
	Id           int     `validate:"min=0"`
	Lat          float64 `validate:"latitude,required"`
	Lon          float64 `validate:"longitude,required"`
	AddressLine1 string  `validate:"required,max=40"`
	AddressLine2 string  `validate:"max=40"`
	City         string  `validate:"required,max=40"`
	Country      string  `validate:"required,max=20"`
}

type AddAddressRequest struct {
	Lat          float64 `validate:"latitude,required"`
	Lon          float64 `validate:"longitude,required"`
	AddressLine1 string  `validate:"required,max=40"`
	AddressLine2 string  `validate:"max=40"`
	City         string  `validate:"required,max=40"`
	Country      string  `validate:"required,max=20"`
}

type Address struct {
	Id           int
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
	case http.MethodDelete:
		AuthenticationMiddleware(c.DeleteAddressHandler, c.secretKey)(w, r)
	}
}

func (c *CustomerAddressServer) DeleteAddressHandler(w http.ResponseWriter, r *http.Request) {
	var deleteAddressRequest DeleteAddressRequest
	err := ValidateBody(r.Body, &deleteAddressRequest)
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

	address, err := c.addressStore.GetAddressById(deleteAddressRequest.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrMissingAddress.Error()})
		return
	}

	if address.CustomerId != customerId {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrUnathorizedAction.Error()})
		return
	}

	c.addressStore.DeleteAddressById(deleteAddressRequest.Id)
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
		Id:           address.Id,
		Lat:          address.Lat,
		Lon:          address.Lon,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		Country:      address.Country,
	}

	return getAddressResponse
}
