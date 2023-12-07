package handlers

import (
	"net/http"

	"github.com/VitoNaychev/food-app/reqbuilder"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

func NewDeleteMenuItemRequest(jwt string, body DeleteMenuItemRequest) *http.Request {
	request := reqbuilder.NewRequestWithBody[DeleteMenuItemRequest](
		http.MethodDelete, "/restaurant/menu/", body)
	request.Header.Add("Token", jwt)

	return request
}

func NewUpdateMenuItemRequest(jwt string, menuItem models.MenuItem) *http.Request {
	updateMenuItemRequest := MenuItemToUpdateMenuItemRequest(menuItem)

	request := reqbuilder.NewRequestWithBody[UpdateMenuItemRequest](
		http.MethodPut, "/restaurant/menu/", updateMenuItemRequest)
	request.Header.Add("Token", jwt)

	return request
}

func NewCreateMenuItemRequest(jwt string, menuItem models.MenuItem) *http.Request {
	createMenuItemRequest := MenuItemToCreateMenuItemRequest(menuItem)

	request := reqbuilder.NewRequestWithBody[CreateMenuItemRequest](
		http.MethodPost, "/restaurant/menu/", createMenuItemRequest)
	request.Header.Add("Token", jwt)

	return request
}

func NewGetMenuRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/restaurant/menu/all/", nil)
	request.Header.Add("Token", jwt)

	return request
}
