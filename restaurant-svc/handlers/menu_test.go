package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/reqbuilder"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	td "github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/restaurant-svc/testutil"
	"github.com/VitoNaychev/food-app/restaurant-svc/testutil/tabletests"
	"github.com/VitoNaychev/food-app/validation"
)

type StubMenuStore struct {
	menus []models.MenuItem
}

func (m *StubMenuStore) GetMenuByRestaurantID(restaurantID int) ([]models.MenuItem, error) {
	menu := []models.MenuItem{}

	for _, item := range m.menus {
		if item.RestaurantID == restaurantID {
			menu = append(menu, item)
		}
	}

	return menu, nil
}

func TestMenuEndpointAuthentication(t *testing.T) {
	restaurantStore := &StubRestaurantStore{}

	menuStore := &StubMenuStore{}

	server := handlers.NewMenuServer(testEnv.SecretKey, menuStore, restaurantStore)
	invalidJWT := "invalidJWT"

	cases := map[string]*http.Request{
		"get menu": NewGetMenuRequest(invalidJWT),
	}

	tabletests.RunAuthenticationTests(t, &server, cases)
}

func TestDeleteMenuItem(t *testing.T) {
	t.Run("deletes menu item on DELETE", func(t *testing.T) {

	})
}

func NewDeleteMenuItemRequest(jwt string, body handlers.DeleteMenuItemRequest) *http.Request {
	request := reqbuilder.NewRequestWithBody[handlers.DeleteMenuItemRequest](
		http.MethodDelete, "/restaurant/menu/", body)
	request.Header.Add("Token", jwt)

	return request
}

func TestUpdateMenuItem(t *testing.T) {
	t.Run("updates menu item on PUT", func(t *testing.T) {

	})
}

func NewUpdateMenuItemRequest(jwt string, body handlers.UpdateMenuItemRequest) *http.Request {
	request := reqbuilder.NewRequestWithBody[handlers.UpdateMenuItemRequest](
		http.MethodPut, "/restaurant/menu/", body)
	request.Header.Add("Token", jwt)

	return request
}

func TestCreateMenuItem(t *testing.T) {
	t.Run("creates menu item on POST", func(t *testing.T) {

	})
}

func NewCreateMenuItemRequest(jwt string, body handlers.CreateMenuItemRequest) *http.Request {
	request := reqbuilder.NewRequestWithBody[handlers.CreateMenuItemRequest](
		http.MethodPut, "/restaurant/menu/", body)
	request.Header.Add("Token", jwt)

	return request
}

func TestGetMenuItem(t *testing.T) {
	t.Run("gets menu item on GET", func(t *testing.T) {

	})
}

func NewGetMenuItemRequest(jwt string, body handlers.GetMenuItemRequest) *http.Request {
	request := reqbuilder.NewRequestWithBody[handlers.GetMenuItemRequest](
		http.MethodGet, "/restaurant/menu/", body)
	request.Header.Add("Token", jwt)

	return request
}

func TestGetMenu(t *testing.T) {
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{td.DominosRestaurant},
	}

	menuStore := &StubMenuStore{
		menus: td.DominosMenu,
	}

	server := handlers.NewMenuServer(testEnv.SecretKey, menuStore, restaurantStore)

	t.Run("gets menu on GET", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)
		request := NewGetMenuRequest(dominosJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]models.MenuItem](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, td.DominosMenu)
	})
}

func NewGetMenuRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/restaurant/menu/all/", nil)
	request.Header.Add("Token", jwt)

	return request
}
