package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/VitoNaychev/food-app/testutil/tabletests"
	"github.com/VitoNaychev/food-app/validation"
)

type StubEventPublisher struct {
	topic string
	event events.InterfaceEvent
}

func (s *StubEventPublisher) Publish(topic string, event events.InterfaceEvent) error {
	s.topic = topic
	s.event = event

	return nil
}

type StubRestaurantStore struct {
	updatedRestaurant   models.Restaurant
	createdRestaurant   models.Restaurant
	deletedRestaurantID int
	restaurants         []models.Restaurant
}

func (s *StubRestaurantStore) DeleteRestaurant(id int) error {
	s.deletedRestaurantID = id
	return nil
}

func (s *StubRestaurantStore) UpdateRestaurant(restaurant *models.Restaurant) error {
	s.updatedRestaurant = *restaurant
	return nil
}

func (s *StubRestaurantStore) CreateRestaurant(restaurant *models.Restaurant) error {
	restaurant.ID = 1
	s.createdRestaurant = *restaurant
	return nil
}

func (s *StubRestaurantStore) GetRestaurantByID(id int) (models.Restaurant, error) {
	for _, restaurant := range s.restaurants {
		if restaurant.ID == id {
			return restaurant, nil
		}
	}

	return models.Restaurant{}, storeerrors.ErrNotFound
}

func (s *StubRestaurantStore) GetRestaurantByEmail(email string) (models.Restaurant, error) {
	for _, restaurant := range s.restaurants {
		if restaurant.Email == email {
			return restaurant, nil
		}
	}

	return models.Restaurant{}, storeerrors.ErrNotFound
}

func TestRestaurantRequestValidation(t *testing.T) {
	store := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.DominosRestaurant},
	}
	server := handlers.NewRestaurantServer(testEnv.SecretKey, testEnv.ExpiresAt, store, nil)

	dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.DominosRestaurant.ID)
	cases := map[string]*http.Request{
		"login restaurant":  handlers.NewLoginRestaurantRequest(models.Restaurant{}),
		"create restaurant": handlers.NewCreateRestaurantRequest(models.Restaurant{}),
		"update restaurant": handlers.NewUpdateRestaurantRequest(dominosJWT, models.Restaurant{}),
	}

	tabletests.RunRequestValidationTests(t, server, cases)
}

func TestRestaurantResponseValidity(t *testing.T) {
	store := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.DominosRestaurant},
	}
	publisher := &StubEventPublisher{}
	server := handlers.NewRestaurantServer(testEnv.SecretKey, testEnv.ExpiresAt, store, publisher)

	dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.DominosRestaurant.ID)
	cases := []tabletests.ResponseValidationTestcase{
		{
			Name:    "get restaurant",
			Request: handlers.NewGetRestaurantRequest(dominosJWT),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[handlers.RestaurantResponse](r)
				return response, err
			},
		},
		{
			Name:    "create restaurant",
			Request: handlers.NewCreateRestaurantRequest(testdata.ShackRestaurant),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[handlers.CreateRestaurantResponse](r)
				return response, err
			},
		},
		{
			Name:    "update restaurant",
			Request: handlers.NewUpdateRestaurantRequest(dominosJWT, testdata.DominosRestaurant),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[handlers.RestaurantResponse](r)
				return response, err
			},
		},
	}

	tabletests.RunResponseValidationTests(t, server, cases)
}

func TestRestaurantEnpointAuthentication(t *testing.T) {
	server := handlers.NewRestaurantServer(testEnv.SecretKey, testEnv.ExpiresAt, nil, nil)

	invalidJWT := "invalidJWT"
	cases := map[string]*http.Request{
		"get restaurant":    handlers.NewGetRestaurantRequest(invalidJWT),
		"update restaurant": handlers.NewUpdateRestaurantRequest(invalidJWT, models.Restaurant{}),
		"delete restaurant": handlers.NewDeleteRestaurantRequest(invalidJWT),
	}

	tabletests.RunAuthenticationTests(t, server, cases)
}

func TestLoginRestaurant(t *testing.T) {
	store := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.ShackRestaurant, testdata.DominosRestaurant},
	}
	server := handlers.NewRestaurantServer(testEnv.SecretKey, testEnv.ExpiresAt, store, nil)

	t.Run("returns JWT on correct credentials", func(t *testing.T) {
		request := handlers.NewLoginRestaurantRequest(testdata.ShackRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		jwtResponse, err := validation.ValidateBody[handlers.JWTResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertJWT(t, jwtResponse.Token, testEnv.SecretKey, testdata.ShackRestaurant.ID)
	})

	t.Run("returns Unauthorized on incorrect email", func(t *testing.T) {
		invalidRestaurant := testdata.ShackRestaurant
		invalidRestaurant.Email = "notshack@gmail.com"

		request := handlers.NewLoginRestaurantRequest(invalidRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidCredentials)
	})

	t.Run("returns Unauthorized on incorrect password", func(t *testing.T) {
		invalidRestaurant := testdata.ShackRestaurant
		invalidRestaurant.Password = "wrongpassword"

		request := handlers.NewLoginRestaurantRequest(invalidRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidCredentials)
	})
}

func TestDeleteRestaurant(t *testing.T) {
	store := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.ShackRestaurant, testdata.DominosRestaurant},
	}
	publisher := &StubEventPublisher{}
	server := handlers.NewRestaurantServer(testEnv.SecretKey, testEnv.ExpiresAt, store, publisher)

	t.Run("deletes restaurant on DELETE", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.DominosRestaurant.ID)

		request := handlers.NewDeleteRestaurantRequest(dominosJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, store.deletedRestaurantID, testdata.DominosRestaurant.ID)
	})

	t.Run("generates a RESTAURANT_DELETED_EVENT on DELETE", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.DominosRestaurant.ID)

		request := handlers.NewDeleteRestaurantRequest(dominosJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertEqual(t, publisher.topic, events.RESTAURANT_EVENTS_TOPIC)

		got := publisher.event
		want := events.InterfaceEvent{
			EventID:     events.RESTAURANT_DELETED_EVENT_ID,
			AggregateID: testdata.DominosAddress.ID,
			Payload:     events.RestaurantDeletedEvent{ID: testdata.DominosRestaurant.ID},
		}

		testutil.AssertEvent(t, got, want)
	})
}

func TestUpdateRestaurant(t *testing.T) {
	store := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.ShackRestaurant, testdata.DominosRestaurant},
	}
	server := handlers.NewRestaurantServer(testEnv.SecretKey, testEnv.ExpiresAt, store, nil)

	t.Run("updates restaurant on PUT", func(t *testing.T) {
		updatedRestaurant := testdata.DominosRestaurant
		updatedRestaurant.Email = "dominos@gmail.com"

		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.DominosRestaurant.ID)

		request := handlers.NewUpdateRestaurantRequest(dominosJWT, updatedRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, store.updatedRestaurant, updatedRestaurant)
	})
}

func TestGetRestaurant(t *testing.T) {
	store := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.ShackRestaurant},
	}
	server := handlers.NewRestaurantServer(testEnv.SecretKey, testEnv.ExpiresAt, store, nil)

	t.Run("resturns restaurant on GET", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.ShackRestaurant.ID)
		request := handlers.NewGetRestaurantRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := handlers.RestaurantToRestaurantResponse(testdata.ShackRestaurant)

		got, err := validation.ValidateBody[handlers.RestaurantResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})
}

func TestCreateRestaurant(t *testing.T) {
	store := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.DominosRestaurant},
	}
	publisher := &StubEventPublisher{}
	server := handlers.NewRestaurantServer(testEnv.SecretKey, testEnv.ExpiresAt, store, publisher)

	t.Run("creates restaurant on POST", func(t *testing.T) {
		request := handlers.NewCreateRestaurantRequest(testdata.ShackRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, store.createdRestaurant, testdata.ShackRestaurant)
	})

	t.Run("returns JWT on POST", func(t *testing.T) {
		request := handlers.NewCreateRestaurantRequest(testdata.ShackRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		createRestaurantResponse, err := validation.ValidateBody[handlers.CreateRestaurantResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		token := createRestaurantResponse.JWT.Token
		testutil.AssertJWT(t, token, testEnv.SecretKey, testdata.ShackRestaurant.ID)
	})

	t.Run("returns the created restaurant on POST", func(t *testing.T) {
		request := handlers.NewCreateRestaurantRequest(testdata.ShackRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		wantRestaurant := handlers.RestaurantToRestaurantResponse(testdata.ShackRestaurant)

		got, err := validation.ValidateBody[handlers.CreateRestaurantResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got.Restaurant, wantRestaurant)
	})

	t.Run("generates a RESTAURANT_CREATED_EVENT on POST", func(t *testing.T) {
		request := handlers.NewCreateRestaurantRequest(testdata.ShackRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertEqual(t, publisher.topic, events.RESTAURANT_EVENTS_TOPIC)

		got := publisher.event
		want := events.InterfaceEvent{
			EventID:     events.RESTAURANT_CREATED_EVENT_ID,
			AggregateID: testdata.ShackRestaurant.ID,
			Payload:     events.RestaurantCreatedEvent{ID: testdata.ShackAddress.ID},
		}

		testutil.AssertEvent(t, got, want)
	})

	t.Run("returns Bad Request on restaurant with same email", func(t *testing.T) {
		request := handlers.NewCreateRestaurantRequest(testdata.DominosRestaurant)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrExistingRestaurant)
	})
}
