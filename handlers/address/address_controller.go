package address

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/handlers/validation"
	"github.com/VitoNaychev/bt-customer-svc/models/address_store"
)

type UpdateAddressRequest struct {
	Id           int     `validate:"min=0"`
	Lat          float64 `validate:"latitude,required"`
	Lon          float64 `validate:"longitude,required"`
	AddressLine1 string  `validate:"required,max=40"`
	AddressLine2 string  `validate:"max=40"`
	City         string  `validate:"required,max=40"`
	Country      string  `validate:"required,max=20"`
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

func (c *CustomerAddressServer) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	var updateAddressRequest UpdateAddressRequest
	err := validation.ValidateBody(r.Body, &updateAddressRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrInvalidRequestField.Error()})
		return
	}

	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err = c.customerStore.GetCustomerById(customerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingCustomer.Error()})
		return
	}

	address, err := c.addressStore.GetAddressById(updateAddressRequest.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingAddress.Error()})
		return
	}

	if address.CustomerId != customerId {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrUnathorizedAction.Error()})
		return
	}

	address = updateAddressRequestToAddress(updateAddressRequest, customerId)

	c.addressStore.UpdateAddress(address)

}

func updateAddressRequestToAddress(UpdateAddressRequest UpdateAddressRequest, customerId int) address_store.Address {
	address := address_store.Address{
		Id:           UpdateAddressRequest.Id,
		CustomerId:   customerId,
		Lat:          UpdateAddressRequest.Lat,
		Lon:          UpdateAddressRequest.Lon,
		AddressLine1: UpdateAddressRequest.AddressLine1,
		AddressLine2: UpdateAddressRequest.AddressLine2,
		City:         UpdateAddressRequest.City,
		Country:      UpdateAddressRequest.Country,
	}

	return address
}

func (c *CustomerAddressServer) DeleteAddressHandler(w http.ResponseWriter, r *http.Request) {
	var deleteAddressRequest DeleteAddressRequest
	err := validation.ValidateBody(r.Body, &deleteAddressRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: err.Error()})
		return
	}

	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err = c.customerStore.GetCustomerById(customerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingCustomer.Error()})
		return
	}

	address, err := c.addressStore.GetAddressById(deleteAddressRequest.Id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingAddress.Error()})
		return
	}

	if address.CustomerId != customerId {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrUnathorizedAction.Error()})
		return
	}

	c.addressStore.DeleteAddressById(deleteAddressRequest.Id)
}

func (c *CustomerAddressServer) StoreAddressHandler(w http.ResponseWriter, r *http.Request) {
	var addAddressRequest AddAddressRequest
	err := validation.ValidateBody(r.Body, &addAddressRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: err.Error()})
		return
	}

	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err = c.customerStore.GetCustomerById(customerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingCustomer.Error()})
		return
	}

	address := addAddressRequestToAddress(addAddressRequest, customerId)

	c.addressStore.StoreAddress(address)
}

func addAddressRequestToAddress(addAddressRequest AddAddressRequest, customerId int) address_store.Address {
	address := address_store.Address{
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
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingCustomer.Error()})
		return
	}

	addresses, _ := c.addressStore.GetAddressesByCustomerId(customerId)

	getAddressResponse := []GetAddressResponse{}
	for _, address := range addresses {
		getAddressResponse = append(getAddressResponse, addressToGetAddressResponse(address))
	}

	json.NewEncoder(w).Encode(getAddressResponse)
}

func addressToGetAddressResponse(address address_store.Address) GetAddressResponse {
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
