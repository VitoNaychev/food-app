package models

type MenuItem struct {
	ID           int `db:"id"`
	RestaurantID int `db:"restaurant_id"`
	Name         string
	Price        float32
}
