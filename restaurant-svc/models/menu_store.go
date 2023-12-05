package models

type MenuStore interface {
	GetMenuByRestaurantID(resturantID int) ([]MenuItem, error)
	CreateMenuItem(*MenuItem) error
}
