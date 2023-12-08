package handlers

import "github.com/VitoNaychev/food-app/restaurant-svc/models"

type DeleteMenuItemRequest struct {
	ID int `validate:"min=1"`
}

type UpdateMenuItemRequest struct {
	ID      int     `validate:"min=1"             json:"id"`
	Name    string  `validate:"min=2,max=20"      json:"name"`
	Price   float32 `validate:"required,max=1000" json:"price"`
	Details string  `validate:"max=1000"          json:"details"`
}

func MenuItemToUpdateMenuItemRequest(item models.MenuItem) UpdateMenuItemRequest {
	request := UpdateMenuItemRequest{
		ID:      item.ID,
		Name:    item.Name,
		Price:   item.Price,
		Details: item.Details,
	}

	return request
}

func UpdateMenuItemRequestToMenuItem(request UpdateMenuItemRequest, restaurantID int) models.MenuItem {
	menuItem := models.MenuItem{
		ID:           request.ID,
		Name:         request.Name,
		Price:        request.Price,
		Details:      request.Details,
		RestaurantID: restaurantID,
	}

	return menuItem
}

type CreateMenuItemRequest struct {
	Name    string  `validate:"min=2,max=20"      json:"name"`
	Price   float32 `validate:"required,max=1000" json:"price"`
	Details string  `validate:"max=1000"          json:"details"`
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
