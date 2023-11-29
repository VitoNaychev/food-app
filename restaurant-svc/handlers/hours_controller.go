package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/errorresponse"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type HoursServer struct {
	secretKey       []byte
	expiresAt       time.Duration
	hoursStore      models.HoursStore
	restaurantStore models.RestaurantStore
}

func NewHoursServer(secretKey []byte, expiresAt time.Duration,
	hoursStore models.HoursStore, restaurantStore models.RestaurantStore) HoursServer {

	hoursServer := HoursServer{
		secretKey:       secretKey,
		expiresAt:       expiresAt,
		hoursStore:      hoursStore,
		restaurantStore: restaurantStore,
	}

	return hoursServer
}

func (h *HoursServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		auth.AuthenticationMiddleware(h.getHours, h.secretKey)(w, r)
	case http.MethodPost:
		auth.AuthenticationMiddleware(h.createHours, h.secretKey)(w, r)
	}
}

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
