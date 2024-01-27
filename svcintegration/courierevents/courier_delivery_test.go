package courierevents

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/courier-svc/handlers"
	"github.com/VitoNaychev/food-app/courier-svc/testdata"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/svcintegration/services"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestCourierDomainEvents(t *testing.T) {
	_, brokersAddrs := integrationutil.SetupKafkaContainer(t)

	courierKeys := []string{"SECRET", "EXPIRES_AT"}
	courierEnv, err := appenv.LoadEnviornment("../../courier-svc/test.env", courierKeys)
	if err != nil {
		t.Fatalf("Failed to load courier env: %v", err)
	}
	courierEnv.KafkaBrokers = brokersAddrs

	deliveryKeys := []string{}
	deliveryEnv, err := appenv.LoadEnviornment("../../delivery-svc/test.env", deliveryKeys)
	if err != nil {
		t.Fatalf("Failed to load delivery env: %v", err)
	}
	deliveryEnv.KafkaBrokers = brokersAddrs

	courierService := services.SetupCourierService(t, courierEnv, ":5050")
	courierService.Run()
	defer courierService.Stop()

	deliveryService := services.SetupDeliveryService(t, deliveryEnv, ":6060")
	deliveryService.Run()
	defer deliveryService.Stop()

	michaelJWT := ""

	t.Run("delivery-svc replicates new courier", func(t *testing.T) {
		request := handlers.NewCreateCourierRequest(testdata.MichaelCourier)
		response := httptest.NewRecorder()

		courierService.Server.Handler.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		want := models.Courier{
			ID:   testdata.MichaelCourier.ID,
			Name: testdata.MichaelCourier.FirstName,
		}
		got, err := deliveryService.CourierStore.GetCourierByID(testdata.MichaelCourier.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)

		michaelJWT = getJWTFromResponseBody(response.Body)
	})

	t.Run("delivery-svc deletes courier", func(t *testing.T) {
		request := handlers.NewDeleteCourierRequest(michaelJWT)
		response := httptest.NewRecorder()

		courierService.Server.Handler.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		time.Sleep(time.Second)

		want := models.Courier{
			ID:   testdata.MichaelCourier.ID,
			Name: testdata.MichaelCourier.FirstName,
		}
		got, err := deliveryService.CourierStore.GetCourierByID(testdata.MichaelCourier.ID)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, got, want)
	})
}

func getJWTFromResponseBody(body io.Reader) string {
	var createCourierResponse handlers.CreateCourierResponse
	json.NewDecoder(body).Decode(&createCourierResponse)

	return createCourierResponse.JWT.Token
}
