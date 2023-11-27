package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/restaurant-svc/testutil"
	"github.com/VitoNaychev/food-app/validation"
)

type StubRestaurantStore struct {
	CreatedRestaurant models.Restaurant
}

func (s *StubRestaurantStore) CreateRestaurant(restaurant *models.Restaurant) error {
	restaurant.ID = 1
	s.CreatedRestaurant = *restaurant
	return nil
}

func TestCreateRestaurant(t *testing.T) {
	store := &StubRestaurantStore{}
	server := handlers.RestaurantServer{
		SecretKey: testEnv.SecretKey,
		ExpiresAt: testEnv.ExpiresAt,
		Store:     store,
	}

	t.Run("creates restaurant on POST", func(t *testing.T) {
		request := NewCreateRestaurantRequest(testdata.Restaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, store.CreatedRestaurant, testdata.Restaurant)

	})

	t.Run("returns JWT on POST", func(t *testing.T) {
		request := NewCreateRestaurantRequest(testdata.Restaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		jwtResponse, err := validation.ValidateBody[handlers.JWTResponse](response.Body)
		if err != nil {
			t.Fatalf("invalid response body: %v", err)
		}

		testutil.AssertJWT(t, jwtResponse.Token, testEnv.SecretKey, testdata.Restaurant.ID)
	})
}

func NewCreateRestaurantRequest(restaurant models.Restaurant) *http.Request {
	createRestaurantRequest := handlers.RestaurantToCreateRestaurantRequest(restaurant)
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(createRestaurantRequest)

	request, _ := http.NewRequest(http.MethodPost, "/restaurant/", body)
	return request
}
