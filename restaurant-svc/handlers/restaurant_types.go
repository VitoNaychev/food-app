package handlers

import "github.com/VitoNaychev/food-app/restaurant-svc/models"

type UpdateRestaurantRequest struct {
	Name        string `validate:"required"`
	PhoneNumber string `validate:"required,phonenumber"`
	Email       string `validate:"required,email"`
	Password    string `valdiate:"required"`
	IBAN        string `validate:"required"`
}

func UpdateRestaurantRequestToRestaurant(request UpdateRestaurantRequest, id int, status models.Status) models.Restaurant {
	restaurant := models.Restaurant{
		ID:          id,
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
		Password:    request.Password,
		IBAN:        request.IBAN,
		Status:      status,
	}

	return restaurant
}

func RestaurantToUpdateRestaurantRequest(restaurant models.Restaurant) UpdateRestaurantRequest {
	updateRestaurantRequest := UpdateRestaurantRequest{
		Name:        restaurant.Name,
		PhoneNumber: restaurant.PhoneNumber,
		Email:       restaurant.Email,
		Password:    restaurant.Password,
		IBAN:        restaurant.IBAN,
	}

	return updateRestaurantRequest
}

type JWTResponse struct {
	Token string
}

type CreateRestaurantResponse struct {
	JWT        JWTResponse
	Restaurant RestaurantResponse
}

type RestaurantResponse struct {
	ID          int
	Name        string
	PhoneNumber string
	Email       string
	IBAN        string
}

func RestaurantToRestaurantResponse(restaurant models.Restaurant) RestaurantResponse {
	restaurantResponse := RestaurantResponse{
		ID:          restaurant.ID,
		Name:        restaurant.Name,
		PhoneNumber: restaurant.PhoneNumber,
		Email:       restaurant.Email,
		IBAN:        restaurant.IBAN,
	}

	return restaurantResponse
}

type CreateRestaurantRequest struct {
	Name        string `validate:"required"`
	PhoneNumber string `validate:"required,phonenumber"`
	Email       string `validate:"required,email"`
	Password    string `valdiate:"required"`
	IBAN        string `validate:"required"`
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
