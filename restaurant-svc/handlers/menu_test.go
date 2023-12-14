package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	td "github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/VitoNaychev/food-app/testutil/tabletests"
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

	return models.MenuItem{}, storeerrors.ErrNotFound
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
		"get menu":         handlers.NewGetMenuRequest(invalidJWT),
		"create menu item": handlers.NewCreateMenuItemRequest(invalidJWT, models.MenuItem{}),
		"udpate menu item": handlers.NewUpdateMenuItemRequest(invalidJWT, models.MenuItem{}),
		"delete menu item": handlers.NewDeleteMenuItemRequest(invalidJWT, handlers.DeleteMenuItemRequest{}),
	}

	tabletests.RunAuthenticationTests(t, server, cases)
}

func TestMenuRequestValdiation(t *testing.T) {
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{td.DominosRestaurant},
	}

	menuStore := &StubMenuStore{}

	server := handlers.NewMenuServer(testEnv.SecretKey, menuStore, restaurantStore)
	dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)

	cases := map[string]*http.Request{
		"create menu item": handlers.NewCreateMenuItemRequest(dominosJWT, models.MenuItem{}),
		"udpate menu item": handlers.NewUpdateMenuItemRequest(dominosJWT, models.MenuItem{}),
		"delete menu item": handlers.NewDeleteMenuItemRequest(dominosJWT, handlers.DeleteMenuItemRequest{}),
	}

	tabletests.RunRequestValidationTests(t, server, cases)
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

		request := handlers.NewDeleteMenuItemRequest(dominosJWT, handlers.DeleteMenuItemRequest{deleteMenuItemID})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertEqual(t, menuStore.deleteMenuItemID, deleteMenuItemID)
	})

	t.Run("returns Not Found on attempt to delete menu item that doesn't exist", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)
		deleteMenuItemID := 10

		request := handlers.NewDeleteMenuItemRequest(dominosJWT, handlers.DeleteMenuItemRequest{deleteMenuItemID})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingMenuItem)
	})

	t.Run("returns Unauthorized on attempt to delete menu item of another restaurant", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)
		deleteMenuItemID := td.ForeignMenuItem.ID

		request := handlers.NewDeleteMenuItemRequest(dominosJWT, handlers.DeleteMenuItemRequest{deleteMenuItemID})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})

	t.Run("returns Bad Request on restaurant with not VALID state", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.ShackRestaurant.ID)
		request := handlers.NewDeleteMenuItemRequest(shackJWT, handlers.DeleteMenuItemRequest{ID: 1})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidRestaurant)
	})
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

		request := handlers.NewUpdateMenuItemRequest(dominosJWT, menuItem)
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

		request := handlers.NewUpdateMenuItemRequest(dominosJWT, menuItem)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrMissingMenuItem)
	})

	t.Run("returns Unauthorized on attempt to update menu item of another restaurant", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)
		menuItem := td.ForeignMenuItem
		menuItem.Name = "Master Burger Pizza"

		request := handlers.NewUpdateMenuItemRequest(dominosJWT, menuItem)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})

	t.Run("returns Bad Request on restaurant with not VALID state", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.ShackRestaurant.ID)
		request := handlers.NewUpdateMenuItemRequest(shackJWT, models.MenuItem{})
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidRestaurant)
	})
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

		request := handlers.NewCreateMenuItemRequest(dominosJWT, menuItem)
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
		request := handlers.NewCreateMenuItemRequest(shackJWT, menuItem)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidRestaurant)
	})
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
		request := handlers.NewGetMenuRequest(dominosJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]models.MenuItem](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, td.DominosMenu)
	})

	t.Run("returns Bad Request on restaurant with not VALID state", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.ShackRestaurant.ID)
		request := handlers.NewGetMenuRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidRestaurant)
	})
}
