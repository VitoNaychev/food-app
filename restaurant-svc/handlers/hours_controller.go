package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/errorresponse"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/validation"
)

func (h *HoursServer) updateHours(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	restaurant, err := h.restaurantStore.GetRestaurantByID(restaurantID)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	if restaurant.Status&models.HOURS_SET == 0 {
		errorresponse.WriteJSONError(w, http.StatusBadRequest, ErrHoursNotSet)
		return
	}

	updateHoursRequestArr, err := validation.ValidateBody[[]HoursRequest](r.Body)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	currentHoursArr, err := h.hoursStore.GetHoursByRestaurantID(restaurantID)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	for i := range currentHoursArr {
		if currentHoursArr[i].Day != i+1 {
			targetIndex := currentHoursArr[i].Day - 1

			temp := currentHoursArr[i]
			currentHoursArr[i] = currentHoursArr[targetIndex]
			currentHoursArr[targetIndex] = temp

			i--
		}
	}

	for i := range updateHoursRequestArr {
		dayIndex := updateHoursRequestArr[i].Day - 1
		currentHours := currentHoursArr[dayIndex]

		updateHours := HoursRequestToHours(updateHoursRequestArr[dayIndex], restaurantID)
		updateHours.ID = currentHours.ID

		err := h.hoursStore.UpdateHours(&updateHours)
		if err != nil {
			errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
			return
		}
	}
}

func (h *HoursServer) createHours(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	restaurant, _ := h.restaurantStore.GetRestaurantByID(restaurantID)
	if restaurant.Status&models.HOURS_SET != 0 {
		errorresponse.WriteJSONError(w, http.StatusBadRequest, ErrHoursAlreadySet)
		return
	}

	var createHoursRequestArr []HoursRequest
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
		hours := HoursRequestToHours(createHoursRequest, restaurantID)
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
