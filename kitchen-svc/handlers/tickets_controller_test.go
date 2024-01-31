package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/VitoNaychev/food-app/appenv"
	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/kitchen-svc/stubs"
	"github.com/VitoNaychev/food-app/kitchen-svc/testdata"

	"github.com/VitoNaychev/food-app/testutil"
	"github.com/VitoNaychev/food-app/testutil/tabletests"
	"github.com/VitoNaychev/food-app/validation"
)

var env appenv.Enviornment

func TestMain(m *testing.M) {
	keys := []string{"SECRET", "EXPIRES_AT"}

	var err error
	env, err = appenv.LoadEnviornment("../test.env", keys)
	if err != nil {
		testutil.HandleLoadEnviornmentError(err)
		os.Exit(1)
	}

	code := m.Run()
	os.Exit(code)
}

func TestTicketEndpointAuthentication(t *testing.T) {
	server := handlers.NewTicketServer(env.SecretKey,
		&stubs.StubTicketStore{},
		&stubs.StubTicketItemStore{},
		&stubs.StubMenuItemStore{},
		&stubs.StubRestaurantStore{},
		&stubs.StubEventPublisher{})

	invalidJWT := "invalidJWT"
	cases := map[string]*http.Request{
		"get tickets":         handlers.NewGetTicketsRequest(invalidJWT, ""),
		"change ticket state": handlers.NewChangeTicketStateRequest(invalidJWT, 1, models.APPROVE_TICKET),
	}

	tabletests.RunAuthenticationTests(t, server, cases)
}

func TestTicketStateTransisions(t *testing.T) {
	ticketStore := &stubs.StubTicketStore{}
	restaurantStore := &stubs.StubRestaurantStore{
		Restaurants: []models.Restaurant{
			testdata.ShackRestaurant,
		},
	}

	publisher := &stubs.StubEventPublisher{}

	server := handlers.NewTicketServer(env.SecretKey, ticketStore, nil, nil, restaurantStore, publisher)

	shackJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.ShackRestaurant.ID)

	t.Run("changes ticket state to IN_PROGRESS on event BEGIN_PREPARING", func(t *testing.T) {
		ticketStore.SpyTicket = testdata.OpenShackTicket

		readyByStr := "23:59"
		readyByTime, _ := handlers.ParseTimeAndSetDate(readyByStr)

		request := handlers.NewBeginPreparingTicketRequest(shackJWT, testdata.OpenShackTicket.ID, readyByStr)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertEqual(t, ticketStore.SpyTicket.State, models.IN_PROGRESS)
		testutil.AssertEqual(t, ticketStore.SpyTicket.ReadyBy, readyByTime)

		want := handlers.StateTransitionResponse{
			ID:    testdata.OpenShackTicket.ID,
			State: "in_progress",
		}

		var got handlers.StateTransitionResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)

		wantTopic := svcevents.KITCHEN_EVENTS_TOPIC
		wantEvent := events.InterfaceEvent{
			EventID:     svcevents.TICKET_BEGIN_PREPARING_EVENT_ID,
			AggregateID: testdata.OpenShackTicket.ID,
			Payload: svcevents.TicketBeginPreparingEvent{
				ID:      testdata.OpenShackTicket.ID,
				ReadyBy: readyByTime,
			},
		}

		testutil.AssertEqual(t, publisher.SpyTopic, wantTopic)
		testutil.AssertEvent(t, publisher.SpyEvent, wantEvent)
	})

	t.Run("returns Unauthorized on restaurant trying to change ticket that it doesn't own", func(t *testing.T) {
		ticketStore.SpyTicket = testdata.ForeginRestaurantTicket

		request := handlers.NewChangeTicketStateRequest(shackJWT, testdata.ForeginRestaurantTicket.ID, models.BEGIN_PREPARING)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})

	t.Run("returns Bad Request attempt to transition to an unsupported state ", func(t *testing.T) {
		ticketStore.SpyTicket = testdata.InProgressShackTicket

		request := handlers.NewChangeTicketStateRequest(shackJWT, testdata.InProgressShackTicket.ID, models.BEGIN_PREPARING)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnsuportedStateTransition)
	})

	t.Run("changes ticket state to DECLINED on event DECLINE", func(t *testing.T) {
		ticketStore.SpyTicket = testdata.OpenShackTicket

		request := handlers.NewChangeTicketStateRequest(shackJWT, testdata.OpenShackTicket.ID, models.DECLINE_TICKET)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertEqual(t, ticketStore.SpyTicket.State, models.DECLINED)

		want := handlers.StateTransitionResponse{
			ID:    testdata.OpenShackTicket.ID,
			State: "declined",
		}

		var got handlers.StateTransitionResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("changes ticket state to READY_FOR_PICKUP on event FINISH_PREPARING", func(t *testing.T) {
		ticketStore.SpyTicket = testdata.InProgressShackTicket

		request := handlers.NewChangeTicketStateRequest(shackJWT, testdata.InProgressShackTicket.ID, models.FINISH_PREPARING)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)
		testutil.AssertEqual(t, ticketStore.SpyTicket.State, models.READY_FOR_PICKUP)

		want := handlers.StateTransitionResponse{
			ID:    testdata.InProgressShackTicket.ID,
			State: "ready_for_pickup",
		}

		var got handlers.StateTransitionResponse
		json.NewDecoder(response.Body).Decode(&got)

		testutil.AssertEqual(t, got, want)

		wantTopic := svcevents.KITCHEN_EVENTS_TOPIC
		wantEvent := events.InterfaceEvent{
			EventID:     svcevents.TICKET_FINISH_PREPARING_EVENT_ID,
			AggregateID: testdata.InProgressShackTicket.ID,
			Payload: svcevents.TicketFinishPreparingEvent{
				ID: testdata.InProgressShackTicket.ID,
			},
		}

		testutil.AssertEqual(t, publisher.SpyTopic, wantTopic)
		testutil.AssertEvent(t, publisher.SpyEvent, wantEvent)
	})
}

func TestTicketGetters(t *testing.T) {
	ticketStore := &stubs.StubTicketStore{
		Tickets: []models.Ticket{
			testdata.OpenShackTicket, testdata.InProgressShackTicket,
			testdata.ReadyForPickupShackTicket, testdata.CompletedShackTicket,
		},
	}
	ticketItemStore := &stubs.StubTicketItemStore{
		TicketItems: []models.TicketItem{
			testdata.OpenShackTicketItems, testdata.InProgressShackTicketItems,
			testdata.ReadyForPickupShackTicketItems, testdata.CompletedShackTicketItems,
		},
	}
	menuItemStore := &stubs.StubMenuItemStore{
		MenuItems: []models.MenuItem{testdata.ShackMenuItem},
	}
	restaurantStore := &stubs.StubRestaurantStore{
		Restaurants: []models.Restaurant{
			testdata.ShackRestaurant,
		},
	}

	publisher := &stubs.StubEventPublisher{}

	server := handlers.NewTicketServer(env.SecretKey, ticketStore, ticketItemStore, menuItemStore, restaurantStore, publisher)

	shackJWT, _ := auth.GenerateJWT(env.SecretKey, env.ExpiresAt, testdata.ShackRestaurant.ID)

	t.Run("returns open tickets on GET to /tickets?state=open", func(t *testing.T) {
		want := testdata.OpenShackTicketResponse

		request := handlers.NewGetTicketsRequest(shackJWT, "?state=open")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns tickets in progress on GET to /tickets?state=in_progress", func(t *testing.T) {
		want := testdata.InProgressShackTicketResponse

		request := handlers.NewGetTicketsRequest(shackJWT, "?state=in_progress")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns tickets ready for pickup on /tickets?state=ready_for_pickup", func(t *testing.T) {
		want := testdata.ReadyForPickupShackTicketResponse

		request := handlers.NewGetTicketsRequest(shackJWT, "?state=ready_for_pickup")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns completed tickets on /tickets?state=completed", func(t *testing.T) {
		want := testdata.CompletedShackTicketResponse

		request := handlers.NewGetTicketsRequest(shackJWT, "?state=completed")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})
}
