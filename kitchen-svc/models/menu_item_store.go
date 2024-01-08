package models

type MenuItemStore interface {
	CreateMenuItem(*MenuItem) error
	DeleteMenuItem(int) error
	UpdateMenuItem(*MenuItem) error
	DeleteMenuItemWhereRestaurantID(int) error
}
