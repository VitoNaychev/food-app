package models

import "time"

type Hours struct {
	ID           int
	Day          int
	Opening      time.Time
	Closing      time.Time
	RestaurantID int `db:"restaurant_id"`
}
