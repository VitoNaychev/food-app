package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

func createRestaurant(server http.Handler, restaurant models.Restaurant) (string, error) {
	request := handlers.NewCreateRestaurantRequest(restaurant)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		return "", fmt.Errorf("failed to create restaurant, got %v", response.Code)
	}

	var restaurantResponse handlers.CreateRestaurantResponse
	json.NewDecoder(response.Body).Decode(&restaurantResponse)

	return restaurantResponse.JWT.Token, nil
}
