package models

import "github.com/VitoNaychev/food-app/storeerrors"

type InMemoryMenuItemStore struct {
	menuItems []MenuItem
}

func NewInMemoryMenuItemStore() *InMemoryMenuItemStore {
	menuItemStore := InMemoryMenuItemStore{
		menuItems: []MenuItem{},
	}

	return &menuItemStore
}

func (i *InMemoryMenuItemStore) GetMenuItemByID(id int) (MenuItem, error) {
	for _, menuItem := range i.menuItems {
		if menuItem.ID == id {
			return menuItem, nil
		}
	}

	return MenuItem{}, storeerrors.ErrNotFound
}

func (i *InMemoryMenuItemStore) CreateMenuItem(menuItem *MenuItem) error {
	i.menuItems = append(i.menuItems, *menuItem)
	return nil
}

func (i *InMemoryMenuItemStore) UpdateMenuItem(menuItem *MenuItem) error {
	for j, oldMenuItem := range i.menuItems {
		if oldMenuItem.ID == menuItem.ID {
			i.menuItems[j] = *menuItem
			return nil
		}

	}

	return storeerrors.ErrNotFound
}

func (i *InMemoryMenuItemStore) DeleteMenuItem(id int) error {
	for j, menuItem := range i.menuItems {
		if menuItem.ID == id {
			i.menuItems = append(i.menuItems[:j], i.menuItems[j+1:]...)
			return nil
		}

	}

	return storeerrors.ErrNotFound
}

func (i *InMemoryMenuItemStore) DeleteMenuItemWhereRestaurantID(restaurantID int) error {
	newMenuItems := []MenuItem{}
	for _, menuItem := range i.menuItems {
		if menuItem.RestaurantID != restaurantID {
			newMenuItems = append(newMenuItems, menuItem)
		}
	}

	i.menuItems = newMenuItems
	return nil
}
