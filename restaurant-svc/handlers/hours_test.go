package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/restaurant-svc/handlers"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
	"github.com/VitoNaychev/food-app/restaurant-svc/testdata"
	"github.com/VitoNaychev/food-app/restaurant-svc/testutil"
	"github.com/VitoNaychev/food-app/restaurant-svc/testutil/tabletests"
	"github.com/VitoNaychev/food-app/validation"
)

type StubHoursStore struct {
	hours        []models.Hours
	createdHours []models.Hours
	updatedHours []models.Hours
}

func (s *StubHoursStore) CreateHours(hours *models.Hours) error {
	hours.ID = len(s.createdHours) + 1
	s.createdHours = append(s.createdHours, *hours)

	return nil
}

func (s *StubHoursStore) UpdateHours(hours *models.Hours) error {
	s.updatedHours = append(s.updatedHours, *hours)

	return nil
}

func (s *StubHoursStore) GetHoursByRestaurantID(restaurantID int) ([]models.Hours, error) {
	days := []models.Hours{}
	for _, day := range s.hours {
		if day.RestaurantID == restaurantID {
			days = append(days, day)
		}
	}

	return days, nil
}

func TestHoursEndpointAuthentication(t *testing.T) {
	server := handlers.NewHoursServer(testEnv.SecretKey, nil, nil)

	cases := map[string]*http.Request{
		"get hours":    handlers.NewGetHoursRequest(""),
		"create hours": handlers.NewCreateHoursRequest("", nil),
		"update hours": handlers.NewUpdateHoursRequest("", nil),
	}

	tabletests.RunAuthenticationTests(t, server, cases)
}

func TestHoursRequestValidation(t *testing.T) {
	hoursStore := &StubHoursStore{
		hours:        []models.Hours{},
		createdHours: []models.Hours{},
	}
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.ShackRestaurant},
	}
	server := handlers.NewHoursServer(testEnv.SecretKey,
		hoursStore,
		restaurantStore)

	shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.ShackRestaurant.ID)

	cases := map[string]*http.Request{
		"create hours": handlers.NewCreateHoursRequest(shackJWT, []models.Hours{}),
		"update hours": handlers.NewUpdateHoursRequest(shackJWT, []models.Hours{}),
	}

	tabletests.RunRequestValidationTests(t, server, cases)
}

func TestUpdateHours(t *testing.T) {
	hoursStore := &StubHoursStore{
		hours:        testdata.DominosHours,
		updatedHours: []models.Hours{},
	}
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.ShackRestaurant, testdata.DominosRestaurant},
	}
	server := handlers.NewHoursServer(testEnv.SecretKey,
		hoursStore,
		restaurantStore)

	t.Run("updates hours on PUT", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.DominosRestaurant.ID)

		updatedHours := make([]models.Hours, 2)
		copy(updatedHours, testdata.DominosHours[:2])
		updatedHours[0].Opening, _ = time.Parse("15:04", "13:00")
		updatedHours[1].Opening, _ = time.Parse("15:04", "13:00")

		request := handlers.NewUpdateHoursRequest(dominosJWT, updatedHours)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertEqual(t, hoursStore.updatedHours, updatedHours)

		assertHoursResponseBody(t, response.Body, updatedHours)
	})

	t.Run("returns Bad Request on update of a restaurant with HOURS_SET bit off", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.ShackRestaurant.ID)

		updatedHours := make([]models.Hours, 2)
		copy(updatedHours, testdata.ShackHours[:2])
		updatedHours[0].Opening, _ = time.Parse("15:04", "13:00")
		updatedHours[1].Opening, _ = time.Parse("15:04", "13:00")

		request := handlers.NewUpdateHoursRequest(shackJWT, updatedHours)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrHoursNotSet)
	})

	t.Run("returns Bad Request on duplicate days", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.DominosRestaurant.ID)

		updatedHours := make([]models.Hours, 2)
		copy(updatedHours, testdata.DominosHours[:1])
		updatedHours[0].Opening, _ = time.Parse("15:04", "13:00")
		updatedHours[1] = updatedHours[0]

		request := handlers.NewUpdateHoursRequest(dominosJWT, updatedHours)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrDuplicateDays)
	})
}

func TestCreateHours(t *testing.T) {
	hoursStore := &StubHoursStore{
		hours:        []models.Hours{},
		createdHours: []models.Hours{},
	}
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.ShackRestaurant, testdata.DominosRestaurant},
	}
	server := handlers.NewHoursServer(testEnv.SecretKey,
		hoursStore,
		restaurantStore)

	t.Run("returns Bad Request if working hours already set", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.DominosRestaurant.ID)
		request := handlers.NewCreateHoursRequest(dominosJWT, testdata.DominosHours)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrHoursAlreadySet)
	})

	t.Run("returns Bad Request if there is a missing day in request", func(t *testing.T) {
		incompleteHours := testdata.ShackHours[1:6]

		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.ShackRestaurant.ID)
		request := handlers.NewCreateHoursRequest(shackJWT, incompleteHours)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrIncompleteWeek)

		restaurant := restaurantStore.updatedRestaurant
		assertRestaurantStatus(t, restaurant, models.CREATED)
	})

	t.Run("returns Bad Request if there are duplicate days in the request", func(t *testing.T) {
		duplicateHours := make([]models.Hours, len(testdata.ShackHours)+1)
		copy(duplicateHours, testdata.ShackHours)
		duplicateHours[7] = duplicateHours[3]

		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.ShackRestaurant.ID)
		request := handlers.NewCreateHoursRequest(shackJWT, duplicateHours)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrDuplicateDays)

		restaurant := restaurantStore.updatedRestaurant
		assertRestaurantStatus(t, restaurant, models.CREATED)
	})

	t.Run("creates working hours for Shack and sets HOURS_SET bit", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.ShackRestaurant.ID)
		request := handlers.NewCreateHoursRequest(shackJWT, testdata.ShackHours)
		response := httptest.NewRecorder()

		// Previous failing tests may have tampered with
		// createdHours, so reinit it to an empty array
		hoursStore.createdHours = []models.Hours{}

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, hoursStore.createdHours, testdata.ShackHours)

		restaurant := restaurantStore.updatedRestaurant
		assertRestaurantStatus(t, restaurant, models.HOURS_SET|models.CREATED)

		assertHoursResponseBody(t, response.Body, testdata.ShackHours)
	})
}

func assertRestaurantStatus(t testing.TB, restaurant models.Restaurant, status models.Status) {
	t.Helper()

	if restaurant.Status == status {
		switch status {
		case models.HOURS_SET | models.HOURS_SET:
			t.Errorf("didn't set HOUR_SET bit in restaurant state")
		case models.CREATED:
			t.Errorf("set HOUR_SET bit in restaurant state, when shouldn't have")
		}
	}
}

func TestGetHours(t *testing.T) {
	hoursStore := &StubHoursStore{
		hours: append(testdata.DominosHours, testdata.ShackHours...),
	}
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.DominosRestaurant, testdata.ShackRestaurant},
	}
	server := handlers.NewHoursServer(testEnv.SecretKey,
		hoursStore,
		restaurantStore)

	t.Run("returns working hours on Chicken Shack", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.ShackRestaurant.ID)
		request := handlers.NewGetHoursRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		assertHoursResponseBody(t, response.Body, testdata.ShackHours)
	})

	t.Run("returns working hours for Dominos", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.DominosRestaurant.ID)
		request := handlers.NewGetHoursRequest(dominosJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		assertHoursResponseBody(t, response.Body, testdata.DominosHours)
	})
}

func assertHoursResponseBody(t testing.TB, body io.Reader, hours []models.Hours) {
	got, err := validation.ValidateBody[[]handlers.HoursResponse](body)
	testutil.AssertValidResponse(t, err)

	want := handlers.HoursArrToHoursResponseArr(hours)

	testutil.AssertEqual(t, got, want)
}
