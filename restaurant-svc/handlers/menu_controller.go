package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (m *MenuServer) getMenu(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	menu, err := m.menuStore.GetMenuByRestaurantID(restaurantID)
	if err != nil {
		handleInternalServerError(w, err)
	}

	json.NewEncoder(w).Encode(menu)
}
