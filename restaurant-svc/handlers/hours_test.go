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
)

type StubHoursStore struct {
	hours        []models.Hours
	createdHours []models.Hours
}

func (s *StubHoursStore) CreateHours(hours *models.Hours) error {
	hours.ID = len(s.createdHours) + 1
	s.createdHours = append(s.createdHours, *hours)

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
	server := handlers.NewHoursServer(testEnv.SecretKey, testEnv.ExpiresAt, nil, nil)

	cases := map[string]*http.Request{
		"get hours":    NewGetHoursRequest(""),
		"create hours": NewCreateHoursRequest("", nil),
	}

	testutil.RunAuthenticationTests(t, &server, cases)
}

func TestCreateHours(t *testing.T) {
	hoursStore := &StubHoursStore{
		hours:        []models.Hours{},
		createdHours: []models.Hours{},
	}
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.ShackRestaurant},
	}
	server := handlers.NewHoursServer(testEnv.SecretKey,
		testEnv.ExpiresAt,
		hoursStore,
		restaurantStore)

	t.Run("returns Bad Request if working hours already set", func(t *testing.T) {

	})

	t.Run("returns Bad Request if the length of the working hours array is not 7", func(t *testing.T) {

	})

	t.Run("creates working hours for Shack and sets HOURS_SET bit", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.ShackRestaurant.ID)
		request := NewCreateHoursRequest(shackJWT, testdata.ShackHours)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, hoursStore.createdHours, testdata.ShackHours)

		restaurant := restaurantStore.updatedRestaurant
		if restaurant.Status&models.HOURS_SET == 0 {
			t.Errorf("didn't set HOUR_SET bit in restaurant state")
		}
	})
}

func NewCreateHoursRequest(jwt string, hours []models.Hours) *http.Request {
	createHoursRequestArr := []handlers.CreateHoursRequest{}
	for _, hour := range hours {
		createHoursRequestArr = append(createHoursRequestArr, handlers.HoursToCreateHoursRequest(hour))
	}

	request := NewRequestWithBody[[]handlers.CreateHoursRequest](
		http.MethodPost, "/restaurant/hours/", createHoursRequestArr)
	SetRequestJWT(request, jwt)

	return request
}

func NewRequestWithBody[T any](method string, url string, object T) *http.Request {
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(object)

	request, _ := http.NewRequest(method, url, body)
	return request
}

func SetRequestJWT(request *http.Request, jwt string) {
	request.Header.Set("Token", jwt)
}

func TestGetHours(t *testing.T) {
	hoursStore := &StubHoursStore{
		hours: append(testdata.DominosHours, testdata.ShackHours...),
	}
	restaurantStore := &StubRestaurantStore{
		restaurants: []models.Restaurant{testdata.DominosRestaurant, testdata.ShackRestaurant},
	}
	server := handlers.NewHoursServer(testEnv.SecretKey,
		testEnv.ExpiresAt,
		hoursStore,
		restaurantStore)

	t.Run("returns working hours on Chicken Shack", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.ShackRestaurant.ID)
		request := NewGetHoursRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		var got []models.Hours
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, testdata.ShackHours)
	})

	t.Run("returns working hours for Dominos", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.DominosRestaurant.ID)
		request := NewGetHoursRequest(dominosJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		var got []models.Hours
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, testdata.DominosHours)
	})

	t.Run("returns Not Found on missing restaurant", func(t *testing.T) {
		missingJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, 10)
		request := NewGetHoursRequest(missingJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func NewGetHoursRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/restaurant/hours/", nil)
	request.Header.Add("Token", jwt)

	return request
}
