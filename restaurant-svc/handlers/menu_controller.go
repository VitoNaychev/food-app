package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
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
		httperrors.HandleBadRequest(w, err)
		return
	}

	currentMenuItem, err := m.menuStore.GetMenuItemByID(deleteMenuItemRequest.ID)
	if err != nil {
		if errors.Is(err, storeerrors.ErrNotFound) {
			httperrors.HandleNotFound(w, ErrMissingMenuItem)
		} else {
			httperrors.HandleInternalServerError(w, err)
		}
		return
	}

	if currentMenuItem.RestaurantID != restaurantID {
		httperrors.HandleUnauthorized(w, ErrUnathorizedAction)
		return
	}

	err = m.menuStore.DeleteMenuItem(deleteMenuItemRequest.ID)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
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
		httperrors.HandleBadRequest(w, err)
		return
	}

	currentMenuItem, err := m.menuStore.GetMenuItemByID(updateMenuItemRequest.ID)
	if err != nil {
		if errors.Is(err, storeerrors.ErrNotFound) {
			httperrors.HandleNotFound(w, ErrMissingMenuItem)
		} else {
			httperrors.HandleInternalServerError(w, err)
		}
		return
	}

	if currentMenuItem.RestaurantID != restaurantID {
		httperrors.HandleUnauthorized(w, ErrUnathorizedAction)
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
		httperrors.HandleBadRequest(w, err)
		return
	}

	menuItem := CreateMenuItemRequestToMenuItem(createMenuItemRequest, restaurantID)

	err = m.menuStore.CreateMenuItem(&menuItem)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
		return
	}

	json.NewEncoder(w).Encode(menuItem)

	event := events.NewEvent(events.RESTAURANT_CREATED_EVENT_ID, restaurantID, events.NewMenuItemCreatedEvent(menuItem, restaurantID))
	m.publisher.Publish(events.RESTAURANT_EVENTS_TOPIC, event)
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
		httperrors.HandleInternalServerError(w, err)
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
		httperrors.HandleBadRequest(w, err)
	} else {
		httperrors.HandleInternalServerError(w, err)
	}
}
