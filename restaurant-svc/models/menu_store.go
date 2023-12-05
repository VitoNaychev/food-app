package models

type MenuStore interface {
	GetMenuByRestaurantID(resturantID int) ([]MenuItem, error)
}
