package address

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/handlers/validation"
)

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

	address = UpdateAddressRequestToAddress(updateAddressRequest, customerId)

	c.addressStore.UpdateAddress(address)

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

	address := AddAddressRequestToAddress(addAddressRequest, customerId)

	c.addressStore.StoreAddress(address)
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
		getAddressResponse = append(getAddressResponse, AddressToGetAddressResponse(address))
	}

	json.NewEncoder(w).Encode(getAddressResponse)
}
