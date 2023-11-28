package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/restaurant-svc/testutil"
	"github.com/VitoNaychev/food-app/validation"
)

type StubRestaurantStore struct {
	CreatedRestaurant models.Restaurant
	restaurants       []models.Restaurant
}

func (s *StubRestaurantStore) CreateRestaurant(restaurant *models.Restaurant) error {
	restaurant.ID = 1
	s.CreatedRestaurant = *restaurant
	return nil
}

func (s *StubRestaurantStore) GetRestaurantByID(id int) (models.Restaurant, error) {
	for _, restaurant := range s.restaurants {
		if restaurant.ID == id {
			return restaurant, nil
		}
	}

	return models.Restaurant{}, models.ErrNotFound
}

func (s *StubRestaurantStore) GetRestaurantByEmail(email string) (models.Restaurant, error) {
	for _, restaurant := range s.restaurants {
		if restaurant.Email == email {
			return restaurant, nil
		}
	}

	return models.Restaurant{}, models.ErrNotFound
}

type DummyRequest struct {
	S string
}

func TestRestaurantRequestValidation(t *testing.T) {
	server := handlers.RestaurantServer{}

	cases := map[string]*http.Request{
		"create restaurant": NewDummyRequest(http.MethodPost),
	}

	for name, request := range cases {
		t.Run(name, func(t *testing.T) {
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		})
	}
}

func NewDummyRequest(method string) *http.Request {
	dummyRequest := DummyRequest{"validation test"}

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(dummyRequest)

	request, _ := http.NewRequest(method, "/restaurant/", body)
	return request
}

func TestRestaurantResponseValidity(t *testing.T) {

}

func TestRestaurantEnpointAuthentication(t *testing.T) {
	server := handlers.RestaurantServer{}

	invalidJWT := "invalidJWT"
	cases := map[string]*http.Request{
		"get restaurant": NewGetRestaurantRequest(invalidJWT),
	}

	for name, request := range cases {
		t.Run(name, func(t *testing.T) {
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		})
	}
}

func TestMissingRestaurantHandling(t *testing.T) {
	store := &StubRestaurantStore{
		restaurants: []models.Restaurant{},
	}
	server := handlers.RestaurantServer{
		SecretKey: testEnv.SecretKey,
		ExpiresAt: testEnv.ExpiresAt,
		Store:     store,
	}

	missingJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 1)

	cases := map[string]*http.Request{
		"get restaurant": NewGetRestaurantRequest(missingJWT),
	}

	for name, request := range cases {
		t.Run(name, func(t *testing.T) {
			response := httptest.NewRecorder()

			server.ServeHTTP(response, request)

			testutil.AssertStatus(t, response.Code, http.StatusNotFound)
		})
	}
}

func TestGetRestaurant(t *testing.T) {
	store := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.ShackRestaurant},
	}
	server := handlers.RestaurantServer{
		SecretKey: testEnv.SecretKey,
		ExpiresAt: testEnv.ExpiresAt,
		Store:     store,
	}

	t.Run("resturns restaurant on GET", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.ShackRestaurant.ID)
		request := NewGetRestaurantRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := handlers.RestaurantToRestaurantResponse(testdata.ShackRestaurant)

		got, err := validation.ValidateBody[handlers.RestaurantResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})

}

func NewGetRestaurantRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/restaurant/", nil)
	request.Header.Add("Token", jwt)

	return request
}

func TestCreateRestaurant(t *testing.T) {
	store := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.DominosRestaurant},
	}
	server := handlers.RestaurantServer{
		SecretKey: testEnv.SecretKey,
		ExpiresAt: testEnv.ExpiresAt,
		Store:     store,
	}

	t.Run("creates restaurant on POST", func(t *testing.T) {
		request := NewCreateRestaurantRequest(testdata.ShackRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, store.CreatedRestaurant, testdata.ShackRestaurant)
	})

	t.Run("returns JWT on POST", func(t *testing.T) {
		request := NewCreateRestaurantRequest(testdata.ShackRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		createRestaurantResponse, err := validation.ValidateBody[handlers.CreateRestaurantResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		token := createRestaurantResponse.JWT.Token
		testutil.AssertJWT(t, token, testEnv.SecretKey, testdata.ShackRestaurant.ID)
	})

	t.Run("returns the created restaurant on POST", func(t *testing.T) {
		request := NewCreateRestaurantRequest(testdata.ShackRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		wantRestaurant := handlers.RestaurantToRestaurantResponse(testdata.ShackRestaurant)

		got, err := validation.ValidateBody[handlers.CreateRestaurantResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got.Restaurant, wantRestaurant)
	})

	t.Run("returns Bad Request on restaurant with same email", func(t *testing.T) {
		request := NewCreateRestaurantRequest(testdata.DominosRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrExistingRestaurant)
	})
}

func NewCreateRestaurantRequest(restaurant models.Restaurant) *http.Request {
	createRestaurantRequest := handlers.RestaurantToCreateRestaurantRequest(restaurant)
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(createRestaurantRequest)

	request, _ := http.NewRequest(http.MethodPost, "/restaurant/", body)
	return request
}
