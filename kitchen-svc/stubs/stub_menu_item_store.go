package stubs

import (
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
)

type StubMenuItemStore struct {
	MenuItems                []models.MenuItem
	CreatedMenuItem          models.MenuItem
	DeletedMenuItemID        int
	UpdatedMenuItem          models.MenuItem
	DeletedItemsRestaurantID int
}

func (s *StubMenuItemStore) GetMenuItemByID(id int) (models.MenuItem, error) {
	for _, menuItem := range s.MenuItems {
		if menuItem.ID == id {
			return menuItem, nil
		}
	}
	return models.MenuItem{}, storeerrors.ErrNotFound
}

func (s *StubMenuItemStore) DeleteMenuItemWhereRestaurantID(id int) error {
	s.DeletedItemsRestaurantID = id
	return nil
}

func (s *StubMenuItemStore) CreateMenuItem(menuItem *models.MenuItem) error {
	s.CreatedMenuItem = *menuItem
	return nil
}

func (s *StubMenuItemStore) DeleteMenuItem(id int) error {
	s.DeletedMenuItemID = id
	return nil
}

func (s *StubMenuItemStore) UpdateMenuItem(menuItem *models.MenuItem) error {
	s.UpdatedMenuItem = *menuItem
	return nil
}
