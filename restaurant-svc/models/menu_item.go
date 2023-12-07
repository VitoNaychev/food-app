package models

type MenuItem struct {
	ID           int
	Name         string
	Price        float32
	Details      string
	RestaurantID int `db:"restaurant_id"`
}
