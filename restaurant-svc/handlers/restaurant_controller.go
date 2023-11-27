package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type RestaurantServer struct {
	SecretKey []byte
	ExpiresAt time.Duration
	Store     models.RestaurantStore
}

func (s *RestaurantServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var createRestaurantRequest CreateRestaurantRequest
	json.NewDecoder(r.Body).Decode(&createRestaurantRequest)

	restaurant := CreateRestaurantRequestToRestaurant(createRestaurantRequest)
	restaurant.Status = models.CREATION_PENDING

	_ = s.Store.CreateRestaurant(&restaurant)

	jwtToken, _ := auth.GenerateJWT(s.SecretKey, s.ExpiresAt, restaurant.ID)

	jwtResponse := JWTResponse{jwtToken}
	json.NewEncoder(w).Encode(jwtResponse)
}
