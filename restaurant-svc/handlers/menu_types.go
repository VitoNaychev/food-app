package handlers

import "github.com/VitoNaychev/food-app/restaurant-svc/models"

type DeleteMenuItemRequest struct {
	ID int `validate:"min=1"`
}

type UpdateMenuItemRequest struct {
	ID      int     `validate:"min=1"`
	Name    string  `validate:"min=2,max=20"`
	Price   float32 `validate:"required,max=1000"`
	Details string  `validate:"max=1000"`
}

type CreateMenuItemRequest struct {
	Name    string  `validate:"min=2,max=20"`
	Price   float32 `validate:"required,max=1000"`
	Details string  `validate:"max=1000"`
}

func MenuItemToCreateMenuItemRequest(item models.MenuItem) CreateMenuItemRequest {
	request := CreateMenuItemRequest{
		Name:    item.Name,
		Price:   item.Price,
		Details: item.Details,
	}

	return request
}

func CreateMenuItemRequestToMenuItem(request CreateMenuItemRequest, restaurantID int) models.MenuItem {
	menuItem := models.MenuItem{
		Name:         request.Name,
		Price:        request.Price,
		Details:      request.Details,
		RestaurantID: restaurantID,
	}

	return menuItem
}
