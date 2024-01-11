package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/integrationutil"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/kitchen-svc/testdata"
	"github.com/VitoNaychev/food-app/pgconfig"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestTicketControllerIntegration(t *testing.T) {
	config := pgconfig.GetConfigFromEnv(env)
	integrationutil.SetupDatabaseContainer(t, &config, "../sql-scripts/init.sql")

	connStr := config.GetConnectionString()

	restaurantStore, err := models.NewPgRestaurantStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)

	menuItemStore, err := models.NewPgMenuItemStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)

	ticketStore, err := models.NewPgTicketStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)

	ticketItemStore, err := models.NewPgTicketItemStore(context.Background(), connStr)
	testutil.AssertNoErr(t, err)

	initTables(t, restaurantStore, menuItemStore, ticketStore, ticketItemStore)

	shackJWT, _ := auth.GenerateJWT(env.SecretKey, time.Second*10, testdata.ShackRestaurant.ID)

	server := handlers.NewTicketServer(env.SecretKey, ticketStore, ticketItemStore, menuItemStore, restaurantStore)

	t.Run("gets all tickets for a restaurant", func(t *testing.T) {
		request := handlers.NewGetTicketsRequest(shackJWT, "")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := []handlers.GetTicketResponse{
			testdata.OpenShackTicketResponse[0],
			testdata.InProgressShackTicketResponse[0],
			testdata.CompletedShackTicketResponse[0],
		}

		var got []handlers.GetTicketResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("marks ticket as ready for pickup", func(t *testing.T) {
		request := handlers.NewChangeTicketStateRequest(shackJWT, testdata.InProgressShackTicket.ID, models.FINISH_PREPARING)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := handlers.StateTransitionResponse{
			ID:    testdata.InProgressShackTicket.ID,
			State: "ready_for_pickup",
		}
		var got handlers.StateTransitionResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("gets all tickets ready for pickup", func(t *testing.T) {
		request := handlers.NewGetTicketsRequest(shackJWT, "?state=ready_for_pickup")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := testdata.InProgressShackTicketResponse
		want[0].State = "ready_for_pickup"

		var got []handlers.GetTicketResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})
}

func initTables(t testing.TB, restaurantStore *models.PgRestaurantStore, menuItemStore *models.PgMenuItemStore,
	ticketStore *models.PgTicketStore, ticketItemStore *models.PgTicketItemStore) {

	testutil.AssertNoErr(t, restaurantStore.CreateRestaurant(&testdata.ShackRestaurant))
	testutil.AssertNoErr(t, menuItemStore.CreateMenuItem(&testdata.ShackMenuItem))

	testutil.AssertNoErr(t, ticketStore.CreateTicket(&testdata.OpenShackTicket))
	testutil.AssertNoErr(t, ticketStore.CreateTicket(&testdata.InProgressShackTicket))
	testutil.AssertNoErr(t, ticketStore.CreateTicket(&testdata.CompletedShackTicket))

	testutil.AssertNoErr(t, ticketItemStore.CreateTicketItem(&testdata.OpenShackTicketItems))
	testutil.AssertNoErr(t, ticketItemStore.CreateTicketItem(&testdata.InProgressShackTicketItems))
	testutil.AssertNoErr(t, ticketItemStore.CreateTicketItem(&testdata.CompletedShackTicketItems))
}
