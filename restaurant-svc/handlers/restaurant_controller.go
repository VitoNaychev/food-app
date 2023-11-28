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
	"github.com/VitoNaychev/food-app/validation"
)

type RestaurantServer struct {
	SecretKey []byte
	ExpiresAt time.Duration
	Store     models.RestaurantStore
}

func (s *RestaurantServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createRestaurant(w, r)
	case http.MethodGet:
		auth.AuthenticationMiddleware(s.getRestaurant, s.SecretKey)(w, r)
	}
}

func (s *RestaurantServer) getRestaurant(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	restaurant, err := s.Store.GetRestaurantByID(restaurantID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			errorresponse.WriteJSONError(w, http.StatusNotFound, err)
		} else {
			errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
		}
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
	restaurant.Status = models.CREATION_PENDING

	if _, err = s.Store.GetRestaurantByEmail(restaurant.Email); !errors.Is(err, models.ErrNotFound) {
		errorresponse.WriteJSONError(w, http.StatusBadRequest, ErrExistingRestaurant)
		return
	}

	err = s.Store.CreateRestaurant(&restaurant)
	if err != nil {
		errorresponse.WriteJSONError(w, http.StatusInternalServerError, err)
	}

	jwtToken, _ := auth.GenerateJWT(s.SecretKey, s.ExpiresAt, restaurant.ID)

	response := CreateRestaurantResponse{
		JWT:        JWTResponse{jwtToken},
		Restaurant: RestaurantToRestaurantResponse(restaurant),
	}
	json.NewEncoder(w).Encode(response)
}
