package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/errorresponse"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

func (h *HoursServer) createHours(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	if !h.CheckIfRestaurantExists(w, restaurantID) {
		return
	}

	var createHoursRequestArr []CreateHoursRequest
	json.NewDecoder(r.Body).Decode(&createHoursRequestArr)

	var hoursArr []models.Hours
	for _, createHoursRequest := range createHoursRequestArr {
		hours := CreateHoursRequestToHours(createHoursRequest, restaurantID)
		hoursArr = append(hoursArr, hours)

		_ = h.hoursStore.CreateHours(&hours)
	}

	restaurant, _ := h.restaurantStore.GetRestaurantByID(restaurantID)
	restaurant.Status = restaurant.Status | models.HOURS_SET
	_ = h.restaurantStore.UpdateRestaurant(&restaurant)
}

func (h *HoursServer) getHours(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	if !h.CheckIfRestaurantExists(w, restaurantID) {
		return
	}

	hours, _ := h.hoursStore.GetHoursByRestaurantID(restaurantID)

	json.NewEncoder(w).Encode(hours)
}

func (h *HoursServer) CheckIfRestaurantExists(w http.ResponseWriter, restaurantID int) bool {
	_, err := h.restaurantStore.GetRestaurantByID(restaurantID)
	if errors.Is(err, models.ErrNotFound) {
		errorresponse.WriteJSONError(w, http.StatusNotFound, err)
		return false
	} else if err != nil {
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
		return false
	}

	return true
}
