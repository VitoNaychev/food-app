package handlers_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	td "github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/restaurant-svc/testutil"
	"github.com/VitoNaychev/food-app/restaurant-svc/testutil/tabletests"
	"github.com/VitoNaychev/food-app/validation"
)

type StubAddressStore struct {
	updatedAddress models.Address
	createdAddress models.Address
	addresses      []models.Address
}

func (s *StubAddressStore) CreateAddress(address *models.Address) error {
	address.ID = 1
	s.createdAddress = *address

	return nil
}

func (s *StubAddressStore) GetAddressByID(id int) (models.Address, error) {
	for _, address := range s.addresses {
		if address.ID == id {
			return address, nil
		}
	}

	return models.Address{}, models.ErrNotFound
}

func (s *StubAddressStore) GetAddressByRestaurantID(restaurantID int) (models.Address, error) {
	for _, address := range s.addresses {
		if address.RestaurantID == restaurantID {
			return address, nil
		}
	}

	return models.Address{}, models.ErrNotFound
}

func (s *StubAddressStore) UpdateAddress(address *models.Address) error {
	s.updatedAddress = *address
	return nil
}

func TestAddressEndpointAuthentication(t *testing.T) {
	addressStore := &StubAddressStore{}
	restaurantStore := &StubRestaurantStore{}

	server := handlers.NewAddressServer(addressStore, restaurantStore, testEnv.SecretKey)

	invalidJWT := "thisIsAnInvalidJWT"
	cases := map[string]*http.Request{
		"get address authentication":    handlers.NewGetAddressRequest(invalidJWT),
		"create address authentication": handlers.NewCreateAddressRequest(invalidJWT, models.Address{}),
		"update address authentication": handlers.NewUpdateAddressRequest(invalidJWT, models.Address{}),
	}

	tabletests.RunAuthenticationTests(t, server, cases)
}

func TestAddressRequestValidation(t *testing.T) {
	addressStore := &StubAddressStore{}
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.DominosRestaurant},
	}

	server := handlers.NewAddressServer(addressStore, restaurantStore, testEnv.SecretKey)

	dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.DominosRestaurant.ID)
	cases := map[string]*http.Request{
		"create address": tabletests.NewDummyRequest(http.MethodPost, "/restaurant/address/", dominosJWT),
		"update address": tabletests.NewDummyRequest(http.MethodPut, "/restaurant/address/", dominosJWT),
	}

	tabletests.RunRequestValidationTests(t, server, cases)
}

func TestAddressResponseValidity(t *testing.T) {
	addressStore := &StubAddressStore{
		addresses: []models.Address{testdata.DominosAddress},
	}
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.ShackRestaurant, testdata.DominosRestaurant},
	}

	server := handlers.NewAddressServer(addressStore, restaurantStore, testEnv.SecretKey)

	shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.ShackRestaurant.ID)
	dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.DominosAddress.ID)
	cases := []tabletests.ResponseValidationTestcase{
		{
			Name:    "get address",
			Request: handlers.NewGetAddressRequest(dominosJWT),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[models.Address](r)
				return response, err
			},
		},
		{
			Name:    "create address",
			Request: handlers.NewCreateAddressRequest(shackJWT, testdata.ShackAddress),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[models.Address](r)
				return response, err
			},
		},
		{
			Name:    "update address",
			Request: handlers.NewUpdateAddressRequest(dominosJWT, testdata.DominosAddress),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[models.Address](r)
				return response, err
			},
		},
	}

	tabletests.RunResponseValidationTests(t, server, cases)
}

func TestUpdateRestaurantAddress(t *testing.T) {
	addressStore := &StubAddressStore{
		addresses: []models.Address{td.DominosAddress},
	}
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{td.ShackRestaurant, td.DominosRestaurant},
	}

	server := handlers.NewAddressServer(addressStore, restaurantStore, testEnv.SecretKey)

	t.Run("updates address on valid body and credentials", func(t *testing.T) {
		updatedAddress := td.DominosAddress
		updatedAddress.City = "Varna"

		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)

		request := handlers.NewUpdateAddressRequest(dominosJWT, updatedAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, addressStore.updatedAddress, updatedAddress)
	})
}

func TestCreateRestaurantAddress(t *testing.T) {
	addressStore := &StubAddressStore{
		addresses: []models.Address{td.DominosAddress},
	}
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{td.ShackRestaurant, td.DominosRestaurant},
	}

	server := handlers.NewAddressServer(addressStore, restaurantStore, testEnv.SecretKey)

	t.Run("creates Shack address and sets ADDRESS_SET bit in restaurant state", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.ShackRestaurant.ID)

		request := handlers.NewCreateAddressRequest(shackJWT, td.ShackAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, addressStore.createdAddress, td.ShackAddress)

		restaurant := restaurantStore.updatedRestaurant
		if restaurant.Status&models.ADDRESS_SET == 0 {
			t.Errorf("didn't set ADDRESS_SET bit in restaurant state")
		}
	})

	t.Run("returns Bad Request if address for restaurant is already set", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)

		request := handlers.NewCreateAddressRequest(dominosJWT, td.DominosAddress)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrAddressAlreadySet)
	})
}

func TestGetRestaurantAddress(t *testing.T) {
	addressStore := &StubAddressStore{
		addresses: []models.Address{td.ShackAddress, td.DominosAddress},
	}
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{td.ShackRestaurant, td.DominosRestaurant},
	}

	server := handlers.NewAddressServer(addressStore, restaurantStore, testEnv.SecretKey)

	t.Run("returns Chicken Shack's address", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.ShackRestaurant.ID)
		request := handlers.NewGetAddressRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := td.ShackAddress

		var got models.Address
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns Dominos address", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, td.DominosRestaurant.ID)
		request := handlers.NewGetAddressRequest(dominosJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := td.DominosAddress

		var got models.Address
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})
}
