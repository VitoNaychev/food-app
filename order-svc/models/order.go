package models

type Order struct {
	ID              int
	CustomerID      int `db:"customer_id"`
	RestaurantID    int `db:"restaurant_id"`
	Total           float32
	Status          Status
	PickupAddress   int `db:"pickup_address"`
	DeliveryAddress int `db:"delivery_address"`
}
