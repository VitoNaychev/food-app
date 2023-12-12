package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/VitoNaychev/food-app/reqbuilder"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

func NewLoginRestaurantRequest(restaurant models.Restaurant) *http.Request {
	requestBody := LoginRestaurantRequest{restaurant.Email, restaurant.Password}
	request := reqbuilder.NewRequestWithBody[LoginRestaurantRequest](
		http.MethodPost, "/restaurant/login/", requestBody)

	return request
}

func NewDeleteRestaruantRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodDelete, "/restaurant/", nil)
	request.Header.Add("Token", jwt)

	return request
}

func NewUpdateRestaurantRequest(jwt string, restaurant models.Restaurant) *http.Request {
	updateRestaurantRequest := RestaurantToUpdateRestaurantRequest(restaurant)
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(updateRestaurantRequest)

	request, _ := http.NewRequest(http.MethodPut, "/restaurant/", body)
	request.Header.Add("Token", jwt)

	return request
}

func NewCreateRestaurantRequest(restaurant models.Restaurant) *http.Request {
	createRestaurantRequest := RestaurantToCreateRestaurantRequest(restaurant)
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(createRestaurantRequest)

	request, _ := http.NewRequest(http.MethodPost, "/restaurant/", body)
	return request
}

func NewGetRestaurantRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/restaurant/", nil)
	request.Header.Add("Token", jwt)

	return request
}
