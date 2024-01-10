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

func TestKitchenEndpointAuthentication(t *testing.T) {

}

func TestTicketStateTransisions(t *testing.T) {
	ticketStore := &StubTicketStore{}

	server := handlers.NewTicketServer(ticketStore, nil, nil)

	t.Run("changes ticket state to PREPARING on POST to /tickets/accept/", func(t *testing.T) {
		ticketStore.spyTicket = testdata.OpenShackTicket
		ticketRequest := handlers.StateTransitionTicketRequest{ID: testdata.OpenShackTicket.ID}

		request := reqbuilder.NewRequestWithBody(http.MethodPost, "/tickets/accept/", ticketRequest)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, ticketStore.spyTicket.State, models.PREPARING)
	})

	t.Run("returns Unauthorized on restaurant trying to change ticket that it doesn't own", func(t *testing.T) {
		ticketStore.spyTicket = testdata.OpenShackTicket
		ticketRequest := handlers.StateTransitionTicketRequest{ID: testdata.OpenShackTicket.ID}

		request := reqbuilder.NewRequestWithBody(http.MethodPost, "/tickets/accept/", ticketRequest)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "5")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusUnauthorized)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnathorizedAction)
	})

	t.Run("returns Bad Request attempt to transition to an unsupported state ", func(t *testing.T) {
		ticketStore.spyTicket = testdata.InProgressShackTicket
		ticketRequest := handlers.StateTransitionTicketRequest{ID: testdata.InProgressShackTicket.ID}

		request := reqbuilder.NewRequestWithBody(http.MethodPost, "/tickets/accept/", ticketRequest)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusBadRequest)
		testutil.AssertErrorResponse(t, response.Body, handlers.ErrUnsuportedStateTransition)
	})

	t.Run("changes ticket state to DECLINED on POST to /tickets/decline/", func(t *testing.T) {
		ticketStore.spyTicket = testdata.OpenShackTicket
		ticketRequest := handlers.StateTransitionTicketRequest{ID: testdata.OpenShackTicket.ID}

		request := reqbuilder.NewRequestWithBody(http.MethodPost, "/tickets/decline/", ticketRequest)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		testutil.AssertEqual(t, ticketStore.spyTicket.State, models.DECLINED)
	})

	t.Run("changes ticket state to READY_FOR_PICKUP on POST to /tickets/prepared/", func(t *testing.T) {
		ticketStore.spyTicket = testdata.InProgressShackTicket
		ticketRequest := handlers.StateTransitionTicketRequest{ID: testdata.InProgressShackTicket.ID}

		request := reqbuilder.NewRequestWithBody(http.MethodPost, "/tickets/prepared/", ticketRequest)
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

	t.Run("returns open tickets on GET to /tickets/open/", func(t *testing.T) {
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

		request, _ := http.NewRequest(http.MethodGet, "/tickets/open/", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns tickets in progress on GET to /tickets/in_progress/", func(t *testing.T) {
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

		request, _ := http.NewRequest(http.MethodGet, "/tickets/in_progress/", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns tickets ready for pickup on /tickets/ready_for_pickup/", func(t *testing.T) {
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

		request, _ := http.NewRequest(http.MethodGet, "/tickets/ready_for_pickup/", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns completed tickets on /tickets/completed/", func(t *testing.T) {
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

		request, _ := http.NewRequest(http.MethodGet, "/tickets/completed/", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})
}
