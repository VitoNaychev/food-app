package models

import "time"

type Order struct {
	ID              int
	CustomerID      int `db:"customer_id"`
	RestaurantID    int `db:"restaurant_id"`
	Items           []int
	Total           float64
	DeliveryTime    time.Time `db:"delivery_time"`
	Status          Status
	PickupAddress   int `db:"pickup_address"`
	DeliveryAddress int `db:"delivery_address"`
}
