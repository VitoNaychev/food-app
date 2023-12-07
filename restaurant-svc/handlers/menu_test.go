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
	menus            []models.MenuItem
	createdMenuItem  models.MenuItem
	updatedMenuItem  models.MenuItem
	deleteMenuItemID int
}

func (m *StubMenuStore) DeleteMenuItem(id int) error {
	m.deleteMenuItemID = id
	return nil
}

func (m *StubMenuStore) UpdateMenuItem(menuItem *models.MenuItem) error {
	m.updatedMenuItem = *menuItem
	return nil
}

func (m *StubMenuStore) CreateMenuItem(menuItem *models.MenuItem) error {
	menuItem.ID = 1
	m.createdMenuItem = *menuItem
	return nil
}

func (m *StubMenuStore) GetMenuItemByID(id int) (models.MenuItem, error) {
	for _, item := range m.menus {
		if item.ID == id {
			return item, nil
		}
	}

	return models.MenuItem{}, models.ErrNotFound
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
		"get menu":         NewGetMenuRequest(invalidJWT),
		"create menu item": NewCreateMenuItemRequest(invalidJWT, models.MenuItem{}),
		"udpate menu item": NewUpdateMenuItemRequest(invalidJWT, models.MenuItem{}),
		"delete menu item": NewDeleteMenuItemRequest(invalidJWT, handlers.DeleteMenuItemRequest{}),
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
		"udpate menu item": NewUpdateMenuItemRequest(dominosJWT, models.MenuItem{}),
		"delete menu item": NewDeleteMenuItemRequest(dominosJWT, handlers.DeleteMenuItemRequest{}),
	}

	tabletests.RunRequestValidationTests(t, &server, cases)
}

func TestDeleteMenuItem(t *testing.T) {
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{td.ShackRestaurant, td.DominosRestaurant},
	}

	menuStore := &StubMenuStore{
		menus: append(td.DominosMenu, td.ForeignMenuItem),
	}

	server := handlers.NewMenuServer(testEnv.SecretKey, menuStore, restaurantStore)

	t.Run("deletes menu item on DELETE", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)
		deleteMenuItemID := td.DominosMenu[1].ID

		request := NewDeleteMenuItemRequest(dominosJWT, handlers.DeleteMenuItemRequest{deleteMenuItemID})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertEqual(t, menuStore.deleteMenuItemID, deleteMenuItemID)
	})

	t.Run("returns Not Found on attempt to delete menu item that doesn't exist", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)
		deleteMenuItemID := 10

		request := NewDeleteMenuItemRequest(dominosJWT, handlers.DeleteMenuItemRequest{deleteMenuItemID})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingMenuItem)
	})

	t.Run("returns Unauthorized on attempt to delete menu item of another restaurant", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)
		deleteMenuItemID := td.ForeignMenuItem.ID

		request := NewDeleteMenuItemRequest(dominosJWT, handlers.DeleteMenuItemRequest{deleteMenuItemID})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})

	t.Run("returns Bad Request on restaurant with not VALID state", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.ShackRestaurant.ID)
		request := NewDeleteMenuItemRequest(shackJWT, handlers.DeleteMenuItemRequest{ID: 1})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidRestaurant)
	})
}

func NewDeleteMenuItemRequest(jwt string, body handlers.DeleteMenuItemRequest) *http.Request {
	request := reqbuilder.NewRequestWithBody[handlers.DeleteMenuItemRequest](
		http.MethodDelete, "/restaurant/menu/", body)
	request.Header.Add("Token", jwt)

	return request
}

func TestUpdateMenuItem(t *testing.T) {
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{td.ShackRestaurant, td.DominosRestaurant},
	}

	menuStore := &StubMenuStore{
		menus: append(td.DominosMenu, td.ForeignMenuItem),
	}

	server := handlers.NewMenuServer(testEnv.SecretKey, menuStore, restaurantStore)

	t.Run("updates menu item on PUT", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)
		menuItem := td.DominosMenu[0]
		menuItem.Name = "Master Burger Pizza"

		request := NewUpdateMenuItemRequest(dominosJWT, menuItem)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, menuStore.updatedMenuItem, menuItem)

		got, err := validation.ValidateBody[models.MenuItem](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, menuItem)
	})

	t.Run("returns Not Found on attempt to update menu item that doesn't exist", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)
		menuItem := models.MenuItem{
			ID:           10,
			Name:         "New Pizza",
			Price:        19.99,
			Details:      "The new-newcomer bruh",
			RestaurantID: td.DominosRestaurant.ID,
		}

		request := NewUpdateMenuItemRequest(dominosJWT, menuItem)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingMenuItem)
	})

	t.Run("returns Unauthorized on attempt to update menu item of another restaurant", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)
		menuItem := td.ForeignMenuItem
		menuItem.Name = "Master Burger Pizza"

		request := NewUpdateMenuItemRequest(dominosJWT, menuItem)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})

	t.Run("returns Bad Request on restaurant with not VALID state", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.ShackRestaurant.ID)
		request := NewUpdateMenuItemRequest(shackJWT, models.MenuItem{})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidRestaurant)
	})
}

func NewUpdateMenuItemRequest(jwt string, menuItem models.MenuItem) *http.Request {
	updateMenuItemRequest := handlers.MenuItemToUpdateMenuItemRequest(menuItem)

	request := reqbuilder.NewRequestWithBody[handlers.UpdateMenuItemRequest](
		http.MethodPut, "/restaurant/menu/", updateMenuItemRequest)
	request.Header.Add("Token", jwt)

	return request
}

func TestCreateMenuItem(t *testing.T) {
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{td.ShackRestaurant, td.DominosRestaurant},
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

	t.Run("returns Bad Request on restaurant with not VALID state", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.ShackRestaurant.ID)
		menuItem := models.MenuItem{
			Name:    "Duner",
			Price:   8.00,
			Details: "on another level",
		}
		request := NewCreateMenuItemRequest(shackJWT, menuItem)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidRestaurant)
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
		restaurants: []models.Restaurant{td.ShackRestaurant, td.DominosRestaurant},
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

	t.Run("returns Bad Request on restaurant with not VALID state", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.ShackRestaurant.ID)
		request := NewGetMenuRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidRestaurant)
	})
}

func NewGetMenuRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/restaurant/menu/all/", nil)
	request.Header.Add("Token", jwt)

	return request
}
