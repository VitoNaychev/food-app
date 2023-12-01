package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/errorresponse"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

func (h *HoursServer) createHours(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	restaurant, _ := h.restaurantStore.GetRestaurantByID(restaurantID)
	if restaurant.Status&models.HOURS_SET != 0 {
		errorresponse.WriteJSONError(w, http.StatusBadRequest, ErrHoursAlreadySet)
		return
	}

	var createHoursRequestArr []CreateHoursRequest
	json.NewDecoder(r.Body).Decode(&createHoursRequestArr)

	var weekBitMask byte
	for _, createHoursRequest := range createHoursRequestArr {
		dayBitMask := byte(1 << (createHoursRequest.Day - 1))
		if weekBitMask&dayBitMask != 0 {
			errorresponse.WriteJSONError(w, http.StatusBadRequest, ErrDuplicateDays)
			return
		}
		weekBitMask |= dayBitMask
	}

	var completeWeekMask byte = 0b01111111
	if weekBitMask&completeWeekMask != completeWeekMask {
		errorresponse.WriteJSONError(w, http.StatusBadRequest, ErrIncompleteWeek)
		return
	}

	var hoursArr []models.Hours
	for _, createHoursRequest := range createHoursRequestArr {
		hours := CreateHoursRequestToHours(createHoursRequest, restaurantID)
		hoursArr = append(hoursArr, hours)

		_ = h.hoursStore.CreateHours(&hours)
	}

	restaurant.Status = restaurant.Status | models.HOURS_SET
	_ = h.restaurantStore.UpdateRestaurant(&restaurant)
}

func (h *HoursServer) getHours(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	hours, _ := h.hoursStore.GetHoursByRestaurantID(restaurantID)

	json.NewEncoder(w).Encode(hours)
}
