package models

type MenuItemStore interface {
	GetMenuItemByID(id int) (MenuItem, error)
	CreateMenuItem(*MenuItem) error
	DeleteMenuItem(int) error
	UpdateMenuItem(*MenuItem) error
	DeleteMenuItemWhereRestaurantID(int) error
}
