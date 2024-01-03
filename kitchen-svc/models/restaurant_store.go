package models

type RestaurantStore interface {
	DeleteRestaurant(int) error
	CreateRestaurant(*Restaurant) error
	GetRestaurantByID(int) (Restaurant, error)
}
