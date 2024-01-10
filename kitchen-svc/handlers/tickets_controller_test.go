package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/kitchen-svc/testdata"
	"github.com/VitoNaychev/food-app/reqbuilder"
	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/VitoNaychev/food-app/validation"
)

type StubTicketStore struct {
	tickets []models.Ticket

	spyTicket models.Ticket
}

func (s *StubTicketStore) GetTicketByID(id int) (models.Ticket, error) {
	if id != s.spyTicket.ID {
		return models.Ticket{}, storeerrors.ErrNotFound
	}

	return s.spyTicket, nil
}

func (s *StubTicketStore) UpdateTicketState(id int, state models.TicketState) error {
	if id != s.spyTicket.ID {
		return storeerrors.ErrNotFound
	}

	s.spyTicket.State = state
	return nil
}

func (s *StubTicketStore) GetTicketsByRestaurantID(restaurantID int) ([]models.Ticket, error) {
	restaruantTickets := []models.Ticket{}

	for _, ticket := range s.tickets {
		if ticket.RestaurantID == restaurantID {
			restaruantTickets = append(restaruantTickets, ticket)
		}
	}

	return restaruantTickets, nil
}

func (s *StubTicketStore) GetTicketsByRestaurantIDWhereState(restaurantID int, state models.TicketState) ([]models.Ticket, error) {
	restaruantTickets := []models.Ticket{}

	for _, ticket := range s.tickets {
		if ticket.RestaurantID == restaurantID &&
			ticket.State == state {
			restaruantTickets = append(restaruantTickets, ticket)
		}
	}

	return restaruantTickets, nil
}

type StubTicketItemStore struct {
	ticketItems []models.TicketItem
}

func (s *StubTicketItemStore) GetTicketItemsByTicketID(ticketID int) ([]models.TicketItem, error) {
	ticketItems := []models.TicketItem{}

	for _, ticketItem := range s.ticketItems {
		if ticketItem.TicketID == ticketID {
			ticketItems = append(ticketItems, ticketItem)
		}
	}

	return ticketItems, nil
}

// func TestTicketEndpointAuthentication(t *testing.T) {
// 	server := handlers.NewTicketServer(testEnv.SecretKey, testEnv.ExpiresAt, nil)

// 	invalidJWT := "invalidJWT"
// 	cases := map[string]*http.Request{
// 		"get open tickets":    ,
// 		"get tickets in progress": ,
// 		"get tickets ready for pickup": ,
// 		"get completed tickets": ,

// 		"accept ticket": ,
// 		"decline ticket": ,
// 		"finish preparing ticket": ,
// 		"complete ticket": ,
// 	}

// 	tabletests.RunAuthenticationTests(t, server, cases)
// }

func TestTicketStateTransisions(t *testing.T) {
	ticketStore := &StubTicketStore{}

	server := handlers.NewTicketServer(ticketStore, nil, nil)

	t.Run("changes ticket state to IN_PROGRESS on event BEGIN_PREPARING", func(t *testing.T) {
		ticketStore.spyTicket = testdata.OpenShackTicket
		ticketRequest := handlers.StateTransitionTicketRequest{
			ID:    testdata.OpenShackTicket.ID,
			Event: models.BEGIN_PREPARING,
		}

		request := reqbuilder.NewRequestWithBody(http.MethodPost, "/tickets/accept/", ticketRequest)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, ticketStore.spyTicket.State, models.IN_PROGRESS)
	})

	t.Run("returns Unauthorized on restaurant trying to change ticket that it doesn't own", func(t *testing.T) {
		ticketStore.spyTicket = testdata.OpenShackTicket
		ticketRequest := handlers.StateTransitionTicketRequest{
			ID:    testdata.OpenShackTicket.ID,
			Event: models.BEGIN_PREPARING,
		}

		request := reqbuilder.NewRequestWithBody(http.MethodPost, "/tickets/accept/", ticketRequest)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "5")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})

	t.Run("returns Bad Request attempt to transition to an unsupported state ", func(t *testing.T) {
		ticketStore.spyTicket = testdata.InProgressShackTicket
		ticketRequest := handlers.StateTransitionTicketRequest{
			ID:    testdata.InProgressShackTicket.ID,
			Event: models.BEGIN_PREPARING,
		}

		request := reqbuilder.NewRequestWithBody(http.MethodPost, "/tickets/accept/", ticketRequest)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnsuportedStateTransition)
	})

	t.Run("changes ticket state to DECLINED on event DECLINE", func(t *testing.T) {
		ticketStore.spyTicket = testdata.OpenShackTicket
		ticketRequest := handlers.StateTransitionTicketRequest{
			ID:    testdata.OpenShackTicket.ID,
			Event: models.DECLINE_TICKET,
		}

		request := reqbuilder.NewRequestWithBody(http.MethodPost, "/tickets/", ticketRequest)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, ticketStore.spyTicket.State, models.DECLINED)
	})

	t.Run("changes ticket state to READY_FOR_PICKUP on event FINISH_PREPARING", func(t *testing.T) {
		ticketStore.spyTicket = testdata.InProgressShackTicket
		ticketRequest := handlers.StateTransitionTicketRequest{
			ID:    testdata.InProgressShackTicket.ID,
			Event: models.FINISH_PREPARING,
		}

		request := reqbuilder.NewRequestWithBody(http.MethodPost, "/tickets/", ticketRequest)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, ticketStore.spyTicket.State, models.READY_FOR_PICKUP)
	})
}

func TestTicketGetters(t *testing.T) {
	ticketStore := &StubTicketStore{
		tickets: []models.Ticket{
			testdata.OpenShackTicket, testdata.InProgressShackTicket, testdata.ReadyForPickupShackTicket, testdata.CompletedShackTicket,
		},
	}
	ticketItemStore := &StubTicketItemStore{
		ticketItems: []models.TicketItem{
			testdata.OpenShackTicketItems, testdata.InProgressShackTicketItems, testdata.ReadyForPickupShackTicketItems, testdata.CompletedShackTicketItems,
		},
	}
	menuItemStore := &StubMenuItemStore{
		menuItems: []models.MenuItem{testdata.ShackMenuItem},
	}

	server := handlers.NewTicketServer(ticketStore, ticketItemStore, menuItemStore)

	t.Run("returns open tickets on GET to /tickets?state=open", func(t *testing.T) {
		want := []handlers.GetTicketResponse{
			{
				ID:    1,
				Total: 12.50,
				Items: []handlers.GetTicketItemResponse{
					{
						Quantity: 2,
						Name:     "Duner",
					},
				},
			},
		}

		request, _ := http.NewRequest(http.MethodGet, "/tickets?state=open", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns tickets in progress on GET to /tickets?state=in_progress", func(t *testing.T) {
		want := []handlers.GetTicketResponse{
			{
				ID:    2,
				Total: 25.00,
				Items: []handlers.GetTicketItemResponse{
					{
						Quantity: 4,
						Name:     "Duner",
					},
				},
			},
		}

		request, _ := http.NewRequest(http.MethodGet, "/tickets?state=in_progress", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns tickets ready for pickup on /tickets?state=ready_for_pickup", func(t *testing.T) {
		want := []handlers.GetTicketResponse{
			{
				ID:    3,
				Total: 37.50,
				Items: []handlers.GetTicketItemResponse{
					{
						Quantity: 6,
						Name:     "Duner",
					},
				},
			},
		}

		request, _ := http.NewRequest(http.MethodGet, "/tickets?state=ready_for_pickup", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns completed tickets on /tickets?state=completed", func(t *testing.T) {
		want := []handlers.GetTicketResponse{
			{
				ID:    4,
				Total: 6.25,
				Items: []handlers.GetTicketItemResponse{
					{
						Quantity: 1,
						Name:     "Duner",
					},
				},
			},
		}

		request, _ := http.NewRequest(http.MethodGet, "/tickets?state=completed", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})
}
