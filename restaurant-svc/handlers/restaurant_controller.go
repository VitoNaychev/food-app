package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/errorresponse"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/validation"
)

func (s *RestaurantServer) updateRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	updateRestaurantRequest, err := validation.ValidateBody[UpdateRestaurantRequest](r.Body)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	oldRestaurant, err := s.store.GetRestaurantByID(restaurantID)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	newRestaurant := UpdateRestaurantRequestToRestaurant(updateRestaurantRequest, restaurantID, oldRestaurant.Status)
	newRestaurant.ID = restaurantID
	newRestaurant.Status = oldRestaurant.Status

	err = s.store.UpdateRestaurant(&newRestaurant)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	updateRestaurantResponse := RestaurantToRestaurantResponse(newRestaurant)
	json.NewEncoder(w).Encode(updateRestaurantResponse)
}

func (s *RestaurantServer) getRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	restaurant, err := s.store.GetRestaurantByID(restaurantID)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
	}

	getRestaurantResponse := RestaurantToRestaurantResponse(restaurant)
	json.NewEncoder(w).Encode(getRestaurantResponse)
}

func (s *RestaurantServer) createRestaurant(w http.ResponseWriter, r *http.Request) {
	createRestaurantRequest, err := validation.ValidateBody[CreateRestaurantRequest](r.Body)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	restaurant := CreateRestaurantRequestToRestaurant(createRestaurantRequest)
	restaurant.Status = models.CREATED

	if _, err = s.store.GetRestaurantByEmail(restaurant.Email); !errors.Is(err, models.ErrNotFound) {
		errorresponse.WriteJSONError(w, http.StatusBadRequest, ErrExistingRestaurant)
		return
	}

	err = s.store.CreateRestaurant(&restaurant)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	jwtToken, _ := auth.GenerateJWT(s.secretKey, s.expiresAt, restaurant.ID)

	response := CreateRestaurantResponse{
		JWT:        JWTResponse{jwtToken},
		Restaurant: RestaurantToRestaurantResponse(restaurant),
	}
	json.NewEncoder(w).Encode(response)
}
