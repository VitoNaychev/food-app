package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/validation"
)

func (m *MenuServer) getMenu(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	menu, err := m.menuStore.GetMenuByRestaurantID(restaurantID)
	if err != nil {
		handleInternalServerError(w, err)
		return
	}

	json.NewEncoder(w).Encode(menu)
}

func (m *MenuServer) createMenuItem(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

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
