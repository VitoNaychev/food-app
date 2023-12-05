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
	menus           []models.MenuItem
	createdMenuItem models.MenuItem
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

func (m *StubMenuStore) CreateMenuItem(menuItem *models.MenuItem) error {
	menuItem.ID = 1
	m.createdMenuItem = *menuItem
	return nil
}

func TestMenuEndpointAuthentication(t *testing.T) {
	restaurantStore := &StubRestaurantStore{}

	menuStore := &StubMenuStore{}

	server := handlers.NewMenuServer(testEnv.SecretKey, menuStore, restaurantStore)
	invalidJWT := "invalidJWT"

	cases := map[string]*http.Request{
		"get menu":         NewGetMenuRequest(invalidJWT),
		"create menu item": NewCreateMenuItemRequest(invalidJWT, models.MenuItem{}),
	}

	tabletests.RunAuthenticationTests(t, &server, cases)
}

func TestMenuRequestValdiation(t *testing.T) {
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{td.DominosRestaurant},
	}

	menuStore := &StubMenuStore{}

	server := handlers.NewMenuServer(testEnv.SecretKey, menuStore, restaurantStore)
	dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)

	cases := map[string]*http.Request{
		"create menu item": NewCreateMenuItemRequest(dominosJWT, models.MenuItem{}),
	}

	tabletests.RunRequestValidationTests(t, &server, cases)
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
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{td.DominosRestaurant},
	}

	menuStore := &StubMenuStore{
		menus: td.DominosMenu,
	}

	server := handlers.NewMenuServer(testEnv.SecretKey, menuStore, restaurantStore)

	t.Run("creates menu item on POST", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)
		menuItem := models.MenuItem{
			Name:    "New Pizza",
			Price:   19.99,
			Details: "The new-newcomer bruh",
		}

		request := NewCreateMenuItemRequest(dominosJWT, menuItem)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := menuItem
		want.ID = 1
		want.RestaurantID = td.DominosRestaurant.ID

		testutil.AssertEqual(t, menuStore.createdMenuItem, want)

		got, err := validation.ValidateBody[models.MenuItem](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})
}

func NewCreateMenuItemRequest(jwt string, menuItem models.MenuItem) *http.Request {
	createMenuItemRequest := handlers.MenuItemToCreateMenuItemRequest(menuItem)

	request := reqbuilder.NewRequestWithBody[handlers.CreateMenuItemRequest](
		http.MethodPost, "/restaurant/menu/", createMenuItemRequest)
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
