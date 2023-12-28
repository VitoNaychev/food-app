package services

import "github.com/VitoNaychev/food-app/kitchen-svc/domain"

type KitchenServiceInterface interface {
	CreateRestaurant(id int) error
	GetRestaurantByID(id int) (domain.Restaurant, error)
	CreateMenuItem(id int, restaurantID int, name string, price float32) error
}
