package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/validation"
)

func (h *HoursServer) updateHours(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	updateHoursRequestArr, err := validation.ValidateBody[[]HoursRequest](r.Body)
	if err != nil {
		httperrors.HandleBadRequest(w, err)
		return
	}

	restaurant, err := h.restaurantStore.GetRestaurantByID(restaurantID)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
		return
	}

	if !isRestaurantStatusBitSet(restaurant, models.HOURS_SET) {
		httperrors.HandleBadRequest(w, ErrHoursNotSet)
		return
	}

	if err := checkForDuplicateDays(updateHoursRequestArr); err != nil {
		httperrors.HandleBadRequest(w, err)
		return
	}

	currentHoursArr, err := h.hoursStore.GetHoursByRestaurantID(restaurantID)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
		return
	}

	updateHoursArr := []models.Hours{}
	for _, updateHoursReq := range updateHoursRequestArr {
		updateHours := HoursRequestToHours(updateHoursReq, -1)
		updateHoursArr = append(updateHoursArr, updateHours)
	}

	setUpdatedHoursKeys(updateHoursArr, currentHoursArr)

	for _, updateHours := range updateHoursArr {
		err := h.hoursStore.UpdateHours(&updateHours)
		if err != nil {
			httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
			return
		}
	}

	updateHoursResponseArr := HoursArrToHoursResponseArr(updateHoursArr)
	json.NewEncoder(w).Encode(updateHoursResponseArr)
}

func setUpdatedHoursKeys(updateHoursArr []models.Hours, currentHoursArr []models.Hours) {
	sortWorkHoursByDay(currentHoursArr)

	for i := range updateHoursArr {
		dayIndex := updateHoursArr[i].Day - 1
		currentHours := currentHoursArr[dayIndex]

		updateHoursArr[i].ID = currentHours.ID
		updateHoursArr[i].RestaurantID = currentHours.RestaurantID
	}
}

func sortWorkHoursByDay(hours []models.Hours) {
	for i := range hours {
		if hours[i].Day != i+1 {
			targetIndex := hours[i].Day - 1

			temp := hours[i]
			hours[i] = hours[targetIndex]
			hours[targetIndex] = temp

			i--
		}
	}
}

func (h *HoursServer) createHours(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	restaurant, err := h.restaurantStore.GetRestaurantByID(restaurantID)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
		return
	}

	if isRestaurantStatusBitSet(restaurant, models.HOURS_SET) {
		httperrors.HandleBadRequest(w, ErrHoursAlreadySet)
		return
	}

	createHoursRequestArr, err := validation.ValidateBody[[]HoursRequest](r.Body)
	if err != nil {
		httperrors.HandleBadRequest(w, err)
		return
	}

	if err = checkForDuplicateOrMissingDays(createHoursRequestArr); err != nil {
		httperrors.HandleBadRequest(w, err)
		return
	}

	var hoursArr []models.Hours
	for _, createHoursRequest := range createHoursRequestArr {
		hours := HoursRequestToHours(createHoursRequest, restaurantID)

		err = h.hoursStore.CreateHours(&hours)
		if err != nil {
			httperrors.HandleInternalServerError(w, err)
			return
		}

		hoursArr = append(hoursArr, hours)
	}

	restaurant.Status = restaurant.Status | models.HOURS_SET
	err = h.restaurantStore.UpdateRestaurant(&restaurant)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
		return
	}

	createHoursResponseArr := HoursArrToHoursResponseArr(hoursArr)
	json.NewEncoder(w).Encode(createHoursResponseArr)
}

func isRestaurantStatusBitSet(restaurant models.Restaurant, status models.Status) bool {
	if restaurant.Status&status != 0 {
		return true
	}

	return false
}

func checkForDuplicateDays(hoursRequestArr []HoursRequest) error {
	var weekBitMask byte
	for _, hoursRequest := range hoursRequestArr {
		dayBitMask := byte(1 << (hoursRequest.Day - 1))
		if weekBitMask&dayBitMask != 0 {
			return ErrDuplicateDays
		}
		weekBitMask |= dayBitMask
	}

	return nil
}

func checkForDuplicateOrMissingDays(hoursRequestArr []HoursRequest) error {
	var weekBitMask byte
	for _, hoursRequest := range hoursRequestArr {
		dayBitMask := byte(1 << (hoursRequest.Day - 1))
		if weekBitMask&dayBitMask != 0 {
			return ErrDuplicateDays
		}
		weekBitMask |= dayBitMask
	}

	var completeWeekMask byte = 0b01111111
	if weekBitMask&completeWeekMask != completeWeekMask {
		return ErrIncompleteWeek
	}

	return nil
}

func (h *HoursServer) getHours(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	hours, err := h.hoursStore.GetHoursByRestaurantID(restaurantID)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
	}

	hoursResponse := HoursArrToHoursResponseArr(hours)
	json.NewEncoder(w).Encode(hoursResponse)
}
