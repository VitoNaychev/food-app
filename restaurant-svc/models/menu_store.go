package models

type MenuStore interface {
	UpdateMenuItem(*MenuItem) error
	CreateMenuItem(*MenuItem) error
	GetMenuItemByID(id int) (MenuItem, error)
	GetMenuByRestaurantID(resturantID int) ([]MenuItem, error)
}
