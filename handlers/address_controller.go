package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/bt-customer-svc/models"
	"github.com/VitoNaychev/validation"
)

func (c *CustomerAddressServer) updateAddress(w http.ResponseWriter, r *http.Request) {
	updateAddressRequest, err := validation.ValidateBody[UpdateAddressRequest](r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, ErrInvalidRequestField)
		return
	}

	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err = c.customerStore.GetCustomerByID(customerId)
	if err != nil {
		handleAddressStoreError(w, err, ErrCustomerNotFound)
		return
	}

	address, err := c.addressStore.GetAddressByID(updateAddressRequest.Id)
	if err != nil {
		handleAddressStoreError(w, err, ErrMissingAddress)
		return
	}

	if address.CustomerId != customerId {
		writeJSONError(w, http.StatusUnauthorized, ErrUnathorizedAction)
		return
	}

	address = UpdateAddressRequestToAddress(updateAddressRequest, customerId)

	err = c.addressStore.UpdateAddress(&address)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, ErrDatabaseError)
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
		handleAddressStoreError(w, err, ErrCustomerNotFound)
		return
	}

	address, err := c.addressStore.GetAddressByID(deleteAddressRequest.Id)
	if err != nil {
		handleAddressStoreError(w, err, ErrMissingAddress)
		return
	}

	if address.CustomerId != customerId {
		writeJSONError(w, http.StatusUnauthorized, ErrUnathorizedAction)
		return
	}

	err = c.addressStore.DeleteAddress(deleteAddressRequest.Id)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, ErrDatabaseError)
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
		handleAddressStoreError(w, err, ErrCustomerNotFound)
		return
	}

	address := CreateAddressRequestToAddress(createAddressRequest, customerId)

	err = c.addressStore.CreateAddress(&address)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, ErrDatabaseError)
	}

	json.NewEncoder(w).Encode(address)
}

func (c *CustomerAddressServer) getAddress(w http.ResponseWriter, r *http.Request) {
	customerId, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err := c.customerStore.GetCustomerByID(customerId)
	if err != nil {
		handleAddressStoreError(w, err, ErrCustomerNotFound)
		return
	}

	addresses, err := c.addressStore.GetAddressesByCustomerID(customerId)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, ErrDatabaseError)
	}

	getAddressResponse := []GetAddressResponse{}
	for _, address := range addresses {
		getAddressResponse = append(getAddressResponse, AddressToGetAddressResponse(address))
	}

	json.NewEncoder(w).Encode(getAddressResponse)
}

func handleAddressStoreError(w http.ResponseWriter, err error, missingEntityError error) {
	if errors.Is(err, models.ErrNotFound) {
		// wrap models.ErrNotFound in customer handlers error type?
		writeJSONError(w, http.StatusNotFound, missingEntityError)
		return
	} else {
		writeJSONError(w, http.StatusInternalServerError, missingEntityError)
		return
	}
}
