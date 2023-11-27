package models

type RestaurantStore interface {
	CreateRestaurant(*Restaurant) error
}
