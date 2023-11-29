package models

type HoursStore interface {
	CreateHours(hours *Hours) error
	GetHoursByRestaurantID(restaurantID int) ([]Hours, error)
}
