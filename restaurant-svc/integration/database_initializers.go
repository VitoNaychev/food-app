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
		return "", fmt.Errorf("failed to create restaurant, got code %v: %v", response.Code, response.Body)
	}

	var restaurantResponse handlers.CreateRestaurantResponse
	json.NewDecoder(response.Body).Decode(&restaurantResponse)

	return restaurantResponse.JWT.Token, nil
}

func createAddress(server http.Handler, jwt string, address models.Address) error {
	request := handlers.NewCreateAddressRequest(jwt, address)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		return fmt.Errorf("failed to create address, got code %v: %v", response.Code, response.Body)
	}

	var restaurantResponse handlers.CreateRestaurantResponse
	json.NewDecoder(response.Body).Decode(&restaurantResponse)

	return nil
}

func createHours(server http.Handler, jwt string, hours []models.Hours) error {
	request := handlers.NewCreateHoursRequest(jwt, hours)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		return fmt.Errorf("failed to create working hours, got code %v: %v", response.Code, response.Body)
	}

	var restaurantResponse handlers.CreateRestaurantResponse
	json.NewDecoder(response.Body).Decode(&restaurantResponse)

	return nil
}
