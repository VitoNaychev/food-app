package models

type RestaurantStore interface {
	CreateRestaurant(*Restaurant) error
	GetRestaurantByID(id int) (Restaurant, error)
	GetRestaurantByEmail(email string) (Restaurant, error)
}
