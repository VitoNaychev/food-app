package handlers_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/courier-svc/handlers"
	"github.com/VitoNaychev/food-app/courier-svc/models"
	"github.com/VitoNaychev/food-app/courier-svc/testdata"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/VitoNaychev/food-app/testutil/tabletests"
	"github.com/VitoNaychev/food-app/validation"
)

type StubCourierStore struct {
	updatedCourier   models.Courier
	createdCourier   models.Courier
	deletedCourierID int
	couriers         []models.Courier
}

func (s *StubCourierStore) DeleteCourier(id int) error {
	s.deletedCourierID = id
	return nil
}

func (s *StubCourierStore) UpdateCourier(courier *models.Courier) error {
	s.updatedCourier = *courier
	return nil
}

func (s *StubCourierStore) CreateCourier(courier *models.Courier) error {
	courier.ID = 1
	s.createdCourier = *courier
	return nil
}

func (s *StubCourierStore) GetCourierByID(id int) (models.Courier, error) {
	for _, courier := range s.couriers {
		if courier.ID == id {
			return courier, nil
		}
	}

	return models.Courier{}, models.ErrNotFound
}

func (s *StubCourierStore) GetCourierByEmail(email string) (models.Courier, error) {
	for _, courier := range s.couriers {
		if courier.Email == email {
			return courier, nil
		}
	}

	return models.Courier{}, models.ErrNotFound
}

func TestCourierRequestValidation(t *testing.T) {
	store := &StubCourierStore{
		couriers: []models.Courier{testdata.JimCourier},
	}
	server := handlers.NewCourierServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.JimCourier.ID)
	cases := map[string]*http.Request{
		"login courier":  handlers.NewLoginCourierRequest(models.Courier{}),
		"create courier": handlers.NewCreateCourierRequest(models.Courier{}),
		"update courier": handlers.NewUpdateCourierRequest(dominosJWT, models.Courier{}),
	}

	tabletests.RunRequestValidationTests(t, server, cases)
}

func TestCourierResponseValidity(t *testing.T) {
	store := &StubCourierStore{
		couriers: []models.Courier{testdata.JimCourier},
	}
	server := handlers.NewCourierServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.JimCourier.ID)
	cases := []tabletests.ResponseValidationTestcase{
		{
			Name:    "get courier",
			Request: handlers.NewGetCourierRequest(dominosJWT),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[handlers.CourierResponse](r)
				return response, err
			},
		},
		{
			Name:    "create courier",
			Request: handlers.NewCreateCourierRequest(testdata.MichaelCourier),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[handlers.CreateCourierResponse](r)
				return response, err
			},
		},
		{
			Name:    "update courier",
			Request: handlers.NewUpdateCourierRequest(dominosJWT, testdata.JimCourier),
			ValidationFunction: func(r io.Reader) (tabletests.GenericResponse, error) {
				response, err := validation.ValidateBody[handlers.CourierResponse](r)
				return response, err
			},
		},
	}

	tabletests.RunResponseValidationTests(t, server, cases)
}

func TestCourierEnpointAuthentication(t *testing.T) {
	server := handlers.NewCourierServer(testEnv.SecretKey, testEnv.ExpiresAt, nil)

	invalidJWT := "invalidJWT"
	cases := map[string]*http.Request{
		"get courier":    handlers.NewGetCourierRequest(invalidJWT),
		"update courier": handlers.NewUpdateCourierRequest(invalidJWT, models.Courier{}),
		"delete courier": handlers.NewDeleteCourierRequest(invalidJWT),
	}

	tabletests.RunAuthenticationTests(t, server, cases)
}

func TestLoginCourier(t *testing.T) {
	store := &StubCourierStore{
		couriers: []models.Courier{testdata.MichaelCourier, testdata.JimCourier},
	}
	server := handlers.NewCourierServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("returns JWT on correct credentials", func(t *testing.T) {
		request := handlers.NewLoginCourierRequest(testdata.MichaelCourier)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		jwtResponse, err := validation.ValidateBody[handlers.JWTResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertJWT(t, jwtResponse.Token, testEnv.SecretKey, testdata.MichaelCourier.ID)
	})

	t.Run("returns Unauthorized on incorrect email", func(t *testing.T) {
		invalidCourier := testdata.MichaelCourier
		invalidCourier.Email = "dwightschrute@gmail.com"

		request := handlers.NewLoginCourierRequest(invalidCourier)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidCredentials)
	})

	t.Run("returns Unauthorized on incorrect password", func(t *testing.T) {
		invalidCourier := testdata.MichaelCourier
		invalidCourier.Password = "wrongpassword"

		request := handlers.NewLoginCourierRequest(invalidCourier)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrInvalidCredentials)
	})
}

func TestDeleteCourier(t *testing.T) {
	store := &StubCourierStore{
		couriers: []models.Courier{testdata.MichaelCourier, testdata.JimCourier},
	}
	server := handlers.NewCourierServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("deletes courier on DELETE", func(t *testing.T) {
		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.JimCourier.ID)

		request := handlers.NewDeleteCourierRequest(dominosJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, store.deletedCourierID, testdata.JimCourier.ID)
	})
}

func TestUpdateCourier(t *testing.T) {
	store := &StubCourierStore{
		couriers: []models.Courier{testdata.MichaelCourier, testdata.JimCourier},
	}
	server := handlers.NewCourierServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("updates courier on PUT", func(t *testing.T) {
		updatedCourier := testdata.JimCourier
		updatedCourier.Email = "prisonmike@gmail.com"

		dominosJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.JimCourier.ID)

		request := handlers.NewUpdateCourierRequest(dominosJWT, updatedCourier)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, store.updatedCourier, updatedCourier)
	})
}

func TestGetCourier(t *testing.T) {
	store := &StubCourierStore{
		couriers: []models.Courier{testdata.MichaelCourier},
	}
	server := handlers.NewCourierServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("returns courier on GET", func(t *testing.T) {
		shackJWT, _ := auth.GenerateJWT(testEnv.SecretKey, testEnv.ExpiresAt, testdata.MichaelCourier.ID)
		request := handlers.NewGetCourierRequest(shackJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := handlers.CourierToCourierResponse(testdata.MichaelCourier)

		got, err := validation.ValidateBody[handlers.CourierResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})
}

func TestCreateCourier(t *testing.T) {
	store := &StubCourierStore{
		couriers: []models.Courier{testdata.JimCourier},
	}
	server := handlers.NewCourierServer(testEnv.SecretKey, testEnv.ExpiresAt, store)

	t.Run("creates courier on POST", func(t *testing.T) {
		request := handlers.NewCreateCourierRequest(testdata.MichaelCourier)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, store.createdCourier, testdata.MichaelCourier)
	})

	t.Run("returns JWT on POST", func(t *testing.T) {
		request := handlers.NewCreateCourierRequest(testdata.MichaelCourier)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		createCourierResponse, err := validation.ValidateBody[handlers.CreateCourierResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		token := createCourierResponse.JWT.Token
		testutil.AssertJWT(t, token, testEnv.SecretKey, testdata.MichaelCourier.ID)
	})

	t.Run("returns the created courier on POST", func(t *testing.T) {
		request := handlers.NewCreateCourierRequest(testdata.MichaelCourier)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		wantCourier := handlers.CourierToCourierResponse(testdata.MichaelCourier)

		got, err := validation.ValidateBody[handlers.CreateCourierResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got.Courier, wantCourier)
	})

	t.Run("returns Bad Request on courier with same email", func(t *testing.T) {
		request := handlers.NewCreateCourierRequest(testdata.JimCourier)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrExistingCourier)
	})
}
