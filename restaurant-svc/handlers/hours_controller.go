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

	updateHoursRequestArr, err := validation.ValidateBody[[]HoursRequest](r.Body)
	if err != nil {
		handleBadRequest(w, err)
		return
	}

	restaurant, err := h.restaurantStore.GetRestaurantByID(restaurantID)
	if err != nil {
		handleInternalServerError(w, err)
		return
	}

	if !isRestaurantStatusBitSet(restaurant, models.HOURS_SET) {
		handleBadRequest(w, ErrHoursNotSet)
		return
	}

	currentHoursArr, err := h.hoursStore.GetHoursByRestaurantID(restaurantID)
	if err != nil {
		handleInternalServerError(w, err)
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
			errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
			return
		}
	}
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
		handleInternalServerError(w, err)
		return
	}

	if isRestaurantStatusBitSet(restaurant, models.HOURS_SET) {
		handleBadRequest(w, ErrHoursAlreadySet)
		return
	}

	var createHoursRequestArr []HoursRequest
	json.NewDecoder(r.Body).Decode(&createHoursRequestArr)

	if err = checkForDuplicateOrMissingDays(createHoursRequestArr); err != nil {
		handleBadRequest(w, err)
	}

	var hoursArr []models.Hours
	for _, createHoursRequest := range createHoursRequestArr {
		hours := HoursRequestToHours(createHoursRequest, restaurantID)
		hoursArr = append(hoursArr, hours)

		err = h.hoursStore.CreateHours(&hours)
		if err != nil {
			handleInternalServerError(w, err)
		}
	}

	restaurant.Status = restaurant.Status | models.HOURS_SET
	err = h.restaurantStore.UpdateRestaurant(&restaurant)
	if err != nil {
		handleInternalServerError(w, err)
		return
	}
}

func isRestaurantStatusBitSet(restaurant models.Restaurant, status models.Status) bool {
	if restaurant.Status&status != 0 {
		return true
	}

	return false
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

	hours, _ := h.hoursStore.GetHoursByRestaurantID(restaurantID)

	json.NewEncoder(w).Encode(hours)
}

func handleBadRequest(w http.ResponseWriter, err error) {
	errorresponse.WriteJSONError(w, http.StatusBadRequest, err)
}

func handleInternalServerError(w http.ResponseWriter, err error) {
	errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
}
