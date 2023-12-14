package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/VitoNaychev/food-app/validation"
)

func (c *AddressServer) updateAddress(w http.ResponseWriter, r *http.Request) {
	updateAddressRequest, err := validation.ValidateBody[UpdateAddressRequest](r.Body)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	restaurant, err := c.restaurantStore.GetRestaurantByID(restaurantID)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	if restaurant.Status&models.ADDRESS_SET == 0 {
		httperrors.WriteJSONError(w, http.StatusNotFound, ErrAddressNotSet)
		return
	}

	currentAddress, err := c.addressStore.GetAddressByRestaurantID(restaurantID)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
	}

	address := UpdateAddressRequestToAddress(updateAddressRequest, currentAddress.ID, restaurantID)

	err = c.addressStore.UpdateAddress(&address)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
	}

	json.NewEncoder(w).Encode(address)
}

func (c *AddressServer) createAddress(w http.ResponseWriter, r *http.Request) {
	createAddressRequest, err := validation.ValidateBody[CreateAddressRequest](r.Body)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	restaurant, err := c.restaurantStore.GetRestaurantByID(restaurantID)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	if restaurant.Status&models.ADDRESS_SET != 0 {
		httperrors.WriteJSONError(w, http.StatusBadRequest, ErrAddressAlreadySet)
		return
	}

	address := CreateAddressRequestToAddress(createAddressRequest, restaurantID)

	err = c.addressStore.CreateAddress(&address)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	restaurant.Status = restaurant.Status | models.ADDRESS_SET
	err = c.restaurantStore.UpdateRestaurant(&restaurant)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	json.NewEncoder(w).Encode(address)
}

func (c *AddressServer) getAddress(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	_, err := c.restaurantStore.GetRestaurantByID(restaurantID)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	address, err := c.addressStore.GetAddressByID(restaurantID)
	if errors.Is(err, storeerrors.ErrNotFound) {
		httperrors.WriteJSONError(w, http.StatusNotFound, err)
		return
	} else if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	json.NewEncoder(w).Encode(address)
}
