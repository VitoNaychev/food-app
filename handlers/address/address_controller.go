package address

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/models"
	"github.com/VitoNaychev/validation"
)

func (c *CustomerAddressServer) updateAddress(w http.ResponseWriter, r *http.Request) {
	updateAddressRequest, err := validation.ValidateBody[UpdateAddressRequest](r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, handlers.ErrInvalidRequestField)
		return
	}

	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err = c.customerStore.GetCustomerByID(customerId)
	if err != nil {
		handleStoreError(w, err, handlers.ErrCustomerNotFound)
		return
	}

	address, err := c.addressStore.GetAddressByID(updateAddressRequest.Id)
	if err != nil {
		handleStoreError(w, err, handlers.ErrMissingAddress)
		return
	}

	if address.CustomerId != customerId {
		writeJSONError(w, http.StatusUnauthorized, handlers.ErrUnathorizedAction)
		return
	}

	address = UpdateAddressRequestToAddress(updateAddressRequest, customerId)

	err = c.addressStore.UpdateAddress(&address)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, handlers.ErrDatabaseError)
	}

	json.NewEncoder(w).Encode(address)
}

func (c *CustomerAddressServer) deleteAddress(w http.ResponseWriter, r *http.Request) {
	deleteAddressRequest, err := validation.ValidateBody[DeleteAddressRequest](r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err = c.customerStore.GetCustomerByID(customerId)
	if err != nil {
		handleStoreError(w, err, handlers.ErrCustomerNotFound)
		return
	}

	address, err := c.addressStore.GetAddressByID(deleteAddressRequest.Id)
	if err != nil {
		handleStoreError(w, err, handlers.ErrMissingAddress)
		return
	}

	if address.CustomerId != customerId {
		writeJSONError(w, http.StatusUnauthorized, handlers.ErrUnathorizedAction)
		return
	}

	err = c.addressStore.DeleteAddress(deleteAddressRequest.Id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, handlers.ErrDatabaseError)
	}
}

func (c *CustomerAddressServer) createAddress(w http.ResponseWriter, r *http.Request) {
	createAddressRequest, err := validation.ValidateBody[CreateAddressRequest](r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err = c.customerStore.GetCustomerByID(customerId)
	if err != nil {
		handleStoreError(w, err, handlers.ErrCustomerNotFound)
		return
	}

	address := CreateAddressRequestToAddress(createAddressRequest, customerId)

	err = c.addressStore.CreateAddress(&address)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, handlers.ErrDatabaseError)
	}

	json.NewEncoder(w).Encode(address)
}

func (c *CustomerAddressServer) getAddress(w http.ResponseWriter, r *http.Request) {
	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err := c.customerStore.GetCustomerByID(customerId)
	if err != nil {
		handleStoreError(w, err, handlers.ErrCustomerNotFound)
		return
	}

	addresses, err := c.addressStore.GetAddressesByCustomerID(customerId)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, handlers.ErrDatabaseError)
	}

	getAddressResponse := []GetAddressResponse{}
	for _, address := range addresses {
		getAddressResponse = append(getAddressResponse, AddressToGetAddressResponse(address))
	}

	json.NewEncoder(w).Encode(getAddressResponse)
}

func writeJSONError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: err.Error()})
}

func handleStoreError(w http.ResponseWriter, err error, missingEntityError error) {
	if errors.Is(err, models.ErrNotFound) {
		// wrap models.ErrNotFound in customer handlers error type?
		writeJSONError(w, http.StatusNotFound, missingEntityError)
		return
	} else {
		writeJSONError(w, http.StatusInternalServerError, missingEntityError)
		return
	}
}
