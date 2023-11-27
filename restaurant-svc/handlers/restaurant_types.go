package handlers

import "github.com/VitoNaychev/food-app/restaurant-svc/models"

type JWTResponse struct {
	Token string
}

type CreateRestaurantRequest struct {
	Name        string
	PhoneNumber string
	Email       string
	Password    string
	IBAN        string
}

func RestaurantToCreateRestaurantRequest(restaurant models.Restaurant) CreateRestaurantRequest {
	createRestaurantRequest := CreateRestaurantRequest{
		Name:        restaurant.Name,
		PhoneNumber: restaurant.PhoneNumber,
		Email:       restaurant.Email,
		Password:    restaurant.Password,
		IBAN:        restaurant.IBAN,
	}

	return createRestaurantRequest
}

func CreateRestaurantRequestToRestaurant(createRestaurantRequest CreateRestaurantRequest) models.Restaurant {
	restaurant := models.Restaurant{
		Name:        createRestaurantRequest.Name,
		PhoneNumber: createRestaurantRequest.PhoneNumber,
		Email:       createRestaurantRequest.Email,
		Password:    createRestaurantRequest.Password,
		IBAN:        createRestaurantRequest.IBAN,
	}

	return restaurant
}
