package address

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/handlers/validation"
)

func (c *CustomerAddressServer) updateAddress(w http.ResponseWriter, r *http.Request) {
	var updateAddressRequest UpdateAddressRequest
	err := validation.ValidateBody(r.Body, &updateAddressRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrInvalidRequestField.Error()})
		return
	}

	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err = c.customerStore.GetCustomerByID(customerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingCustomer.Error()})
		return
	}

	address, err := c.addressStore.GetAddressByID(updateAddressRequest.Id)
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

	err = c.addressStore.UpdateAddress(&address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrDatabaseError.Error()})
	}

	json.NewEncoder(w).Encode(address)
}

func (c *CustomerAddressServer) deleteAddress(w http.ResponseWriter, r *http.Request) {
	var deleteAddressRequest DeleteAddressRequest
	err := validation.ValidateBody(r.Body, &deleteAddressRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: err.Error()})
		return
	}

	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err = c.customerStore.GetCustomerByID(customerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingCustomer.Error()})
		return
	}

	address, err := c.addressStore.GetAddressByID(deleteAddressRequest.Id)
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

	err = c.addressStore.DeleteAddress(deleteAddressRequest.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrDatabaseError.Error()})
	}
}

func (c *CustomerAddressServer) createAddress(w http.ResponseWriter, r *http.Request) {
	var addAddressRequest AddAddressRequest
	err := validation.ValidateBody(r.Body, &addAddressRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: err.Error()})
		return
	}

	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err = c.customerStore.GetCustomerByID(customerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingCustomer.Error()})
		return
	}

	address := AddAddressRequestToAddress(addAddressRequest, customerId)

	err = c.addressStore.CreateAddress(&address)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrDatabaseError.Error()})
	}

	json.NewEncoder(w).Encode(address)
}

func (c *CustomerAddressServer) getAddress(w http.ResponseWriter, r *http.Request) {
	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err := c.customerStore.GetCustomerByID(customerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingCustomer.Error()})
		return
	}

	addresses, _ := c.addressStore.GetAddressesByCustomerID(customerId)

	getAddressResponse := []GetAddressResponse{}
	for _, address := range addresses {
		getAddressResponse = append(getAddressResponse, AddressToGetAddressResponse(address))
	}

	json.NewEncoder(w).Encode(getAddressResponse)
}
