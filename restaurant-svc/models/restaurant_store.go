package models

type RestaurantStore interface {
	DeleteRestaurant(id int) error
	UpdateRestaurant(*Restaurant) error
	CreateRestaurant(*Restaurant) error
	GetRestaurantByID(id int) (Restaurant, error)
	GetRestaurantByEmail(email string) (Restaurant, error)
}
