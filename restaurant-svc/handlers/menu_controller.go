package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/errorresponse"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/validation"
)

func (m *MenuServer) deleteMenuItem(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	err := isRestaurantValid(restaurantID, m.restaurantStore)
	if err != nil {
		handleRestaurantInvalid(w, err)
		return
	}

	deleteMenuItemRequest, err := validation.ValidateBody[DeleteMenuItemRequest](r.Body)
	if err != nil {
		handleBadRequest(w, err)
		return
	}

	currentMenuItem, err := m.menuStore.GetMenuItemByID(deleteMenuItemRequest.ID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			handleNotFoundError(w, ErrMissingMenuItem)
		} else {
			handleInternalServerError(w, err)
		}
		return
	}

	if currentMenuItem.RestaurantID != restaurantID {
		handleUnauthorizedError(w, ErrUnathorizedAction)
		return
	}

	err = m.menuStore.DeleteMenuItem(deleteMenuItemRequest.ID)
	if err != nil {
		handleInternalServerError(w, err)
	}
}

func (m *MenuServer) updateMenuItem(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	err := isRestaurantValid(restaurantID, m.restaurantStore)
	if err != nil {
		handleRestaurantInvalid(w, err)
		return
	}

	updateMenuItemRequest, err := validation.ValidateBody[UpdateMenuItemRequest](r.Body)
	if err != nil {
		handleBadRequest(w, err)
		return
	}

	currentMenuItem, err := m.menuStore.GetMenuItemByID(updateMenuItemRequest.ID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			handleNotFoundError(w, ErrMissingMenuItem)
		} else {
			handleInternalServerError(w, err)
		}
		return
	}

	if currentMenuItem.RestaurantID != restaurantID {
		handleUnauthorizedError(w, ErrUnathorizedAction)
		return
	}

	updateMenuItem := UpdateMenuItemRequestToMenuItem(updateMenuItemRequest, restaurantID)
	err = m.menuStore.UpdateMenuItem(&updateMenuItem)

	json.NewEncoder(w).Encode(updateMenuItem)
}

func (m *MenuServer) createMenuItem(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	err := isRestaurantValid(restaurantID, m.restaurantStore)
	if err != nil {
		handleRestaurantInvalid(w, err)
		return
	}

	createMenuItemRequest, err := validation.ValidateBody[CreateMenuItemRequest](r.Body)
	if err != nil {
		handleBadRequest(w, err)
		return
	}

	menuItem := CreateMenuItemRequestToMenuItem(createMenuItemRequest, restaurantID)

	err = m.menuStore.CreateMenuItem(&menuItem)
	if err != nil {
		handleInternalServerError(w, err)
		return
	}

	json.NewEncoder(w).Encode(menuItem)
}

func (m *MenuServer) getMenu(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	err := isRestaurantValid(restaurantID, m.restaurantStore)
	if err != nil {
		handleRestaurantInvalid(w, err)
		return
	}

	menu, err := m.menuStore.GetMenuByRestaurantID(restaurantID)
	if err != nil {
		handleInternalServerError(w, err)
		return
	}

	json.NewEncoder(w).Encode(menu)
}

func isRestaurantValid(restaurantID int, store models.RestaurantStore) error {
	restaurant, err := store.GetRestaurantByID(restaurantID)
	if err != nil {
		return err
	}

	if restaurant.Status != models.VALID {
		return ErrInvalidRestaurant
	}

	return nil
}

func handleRestaurantInvalid(w http.ResponseWriter, err error) {
	if errors.Is(err, ErrInvalidRestaurant) {
		handleBadRequest(w, err)
	} else {
		handleInternalServerError(w, err)
	}
}

func handleUnauthorizedError(w http.ResponseWriter, err error) {
	errorresponse.WriteJSONError(w, http.StatusUnauthorized, err)
}

func handleNotFoundError(w http.ResponseWriter, err error) {
	errorresponse.WriteJSONError(w, http.StatusNotFound, err)

}
