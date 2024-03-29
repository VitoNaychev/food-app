package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/VitoNaychev/food-app/validation"
)

func (s *RestaurantServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginRestaurantRequest, err := validation.ValidateBody[LoginRestaurantRequest](r.Body)
	if err != nil {
		httperrors.HandleBadRequest(w, err)
		return
	}

	restaurant, err := s.store.GetRestaurantByEmail(loginRestaurantRequest.Email)
	if err != nil {
		if errors.Is(err, storeerrors.ErrNotFound) {
			httperrors.HandleUnauthorized(w, ErrInvalidCredentials)
			return
		} else {
			httperrors.HandleInternalServerError(w, err)
			return
		}
	}

	if restaurant.Password != loginRestaurantRequest.Password {
		httperrors.HandleUnauthorized(w, ErrInvalidCredentials)
		return
	}

	jwtToken, _ := auth.GenerateJWT(s.secretKey, s.expiresAt, restaurant.ID)
	jwtResponse := JWTResponse{Token: jwtToken}

	json.NewEncoder(w).Encode(jwtResponse)
}

func (s *RestaurantServer) deleteRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	err := s.store.DeleteRestaurant(restaurantID)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
	}

	payload := events.RestaurantDeletedEvent{ID: restaurantID}
	event := events.NewEvent(events.RESTAURANT_DELETED_EVENT_ID, restaurantID, payload)
	s.publisher.Publish(events.RESTAURANT_EVENTS_TOPIC, event)
}

func (s *RestaurantServer) updateRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	updateRestaurantRequest, err := validation.ValidateBody[UpdateRestaurantRequest](r.Body)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	oldRestaurant, err := s.store.GetRestaurantByID(restaurantID)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	newRestaurant := UpdateRestaurantRequestToRestaurant(updateRestaurantRequest, restaurantID, oldRestaurant.Status)
	newRestaurant.ID = restaurantID
	newRestaurant.Status = oldRestaurant.Status

	err = s.store.UpdateRestaurant(&newRestaurant)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	updateRestaurantResponse := RestaurantToRestaurantResponse(newRestaurant)
	json.NewEncoder(w).Encode(updateRestaurantResponse)
}

func (s *RestaurantServer) getRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	restaurant, err := s.store.GetRestaurantByID(restaurantID)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
	}

	getRestaurantResponse := RestaurantToRestaurantResponse(restaurant)
	json.NewEncoder(w).Encode(getRestaurantResponse)
}

func (s *RestaurantServer) createRestaurant(w http.ResponseWriter, r *http.Request) {
	createRestaurantRequest, err := validation.ValidateBody[CreateRestaurantRequest](r.Body)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	restaurant := CreateRestaurantRequestToRestaurant(createRestaurantRequest)
	if _, err = s.store.GetRestaurantByEmail(restaurant.Email); !errors.Is(err, storeerrors.ErrNotFound) {
		httperrors.WriteJSONError(w, http.StatusBadRequest, ErrExistingRestaurant)
		return
	}

	restaurant.Status = models.CREATED
	err = s.store.CreateRestaurant(&restaurant)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	jwtToken, _ := auth.GenerateJWT(s.secretKey, s.expiresAt, restaurant.ID)

	response := CreateRestaurantResponse{
		JWT:        JWTResponse{jwtToken},
		Restaurant: RestaurantToRestaurantResponse(restaurant),
	}
	json.NewEncoder(w).Encode(response)

	payload := events.RestaurantCreatedEvent{ID: restaurant.ID}
	event := events.NewEvent(events.RESTAURANT_CREATED_EVENT_ID, restaurant.ID, payload)
	s.publisher.Publish(events.RESTAURANT_EVENTS_TOPIC, event)
}
