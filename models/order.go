package models

import "time"

type Order struct {
	ID              int
	CustomerID      int
	RestaurantID    int
	Items           []int
	Total           float64
	DeliveryTime    time.Time
	PickupAddress   int
	DeliveryAddress int
	Status          Status
}
