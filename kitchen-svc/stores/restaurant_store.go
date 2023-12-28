package stores

import "github.com/VitoNaychev/food-app/kitchen-svc/domain"

type RestaurantStore interface {
	UpdateRestaurant(*domain.Restaurant) error
	CreateRestaurant(*domain.Restaurant) error
	GetRestaurantByID(int) (domain.Restaurant, error)
}
