package services

type KitchenServiceInterface interface {
	CreateRestaurant(id int) error
	CreateMenuItem(id int, restaurantID int, name string, price float32) error
}
