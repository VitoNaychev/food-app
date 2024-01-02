package models

import "github.com/VitoNaychev/food-app/storeerrors"

type InMemoryMenuStore struct {
	menuItems []MenuItem
}

func NewInMemoryMenuStore() *InMemoryMenuStore {
	return &InMemoryMenuStore{[]MenuItem{}}
}

func (i *InMemoryMenuStore) CreateMenuItem(menuItem *MenuItem) error {
	menuItem.ID = len(i.menuItems) + 1
	i.menuItems = append(i.menuItems, *menuItem)
	return nil
}

func (i *InMemoryMenuStore) DeleteMenuItem(id int) error {
	for j, menuItem := range i.menuItems {
		if menuItem.ID == id {
			i.menuItems = append(i.menuItems[:j], i.menuItems[j+1:]...)
			return nil
		}
	}

	return storeerrors.ErrNotFound
}

func (i *InMemoryMenuStore) GetMenuByRestaurantID(resturantID int) ([]MenuItem, error) {
	menuItems := []MenuItem{}
	for _, menuItem := range i.menuItems {
		if menuItem.RestaurantID == resturantID {
			menuItems = append(menuItems, menuItem)
		}
	}

	return menuItems, nil
}

func (i *InMemoryMenuStore) GetMenuItemByID(id int) (MenuItem, error) {
	for _, menuItem := range i.menuItems {
		if menuItem.ID == id {
			return menuItem, nil
		}
	}

	return MenuItem{}, storeerrors.ErrNotFound
}

func (i *InMemoryMenuStore) UpdateMenuItem(menuItem *MenuItem) error {
	for j, oldMenuItem := range i.menuItems {
		if oldMenuItem.ID == menuItem.ID {
			i.menuItems[j] = *menuItem
			return nil
		}
	}

	return storeerrors.ErrNotFound
}
