package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitoNaychev/food-app/kitchen-svc/handlers"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/kitchen-svc/testdata"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/VitoNaychev/food-app/validation"
)

type StubTicketStore struct {
	tickets []models.Ticket
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

// func TestTicketStateTransisions(t *testing.T) {

// 	server := &handlers.TicketServer{}

// 	t.Run("accepts ticket on POST to /tickets/accept/", func(t *testing.T) {
// 		request, _ := http.NewRequest(http.MethodPost, "/tickets/accept/", nil)
// 		response := httptest.NewRecorder()

// 		server.ServeHTTP(response, request)

// 		testutil.AssertStatus(t, response.Code, http.StatusOK)
// 	})
// }

func TestTicketGetters(t *testing.T) {
	ticketStore := &StubTicketStore{
		tickets: []models.Ticket{
			testdata.OpenShackTicket, testdata.InProgressShackTicket, testdata.ReadyForPickupShackTicket, testdata.PendingShackTicket,
		},
	}
	ticketItemStore := &StubTicketItemStore{
		ticketItems: []models.TicketItem{
			testdata.OpenShackTicketItems, testdata.InProgressShackTicketItems, testdata.ReadyForPickupShackTicketItems, testdata.PendingShackTicketItems,
		},
	}
	menuItemStore := &StubMenuItemStore{
		menuItems: []models.MenuItem{testdata.ShackMenuItem},
	}

	server := handlers.NewTicketServer(ticketStore, ticketItemStore, menuItemStore)

	t.Run("returns pending tickets on GET to /tickets/pending/", func(t *testing.T) {
		want := []handlers.GetTicketResponse{
			{
				ID: 4,
				Items: []handlers.GetTicketItemResponse{
					{
						Quantity: 1,
						Name:     "Duner",
					},
				},
			},
		}

		request, _ := http.NewRequest(http.MethodGet, "/tickets/pending/", nil)
		response := httptest.NewRecorder()

		request.Header.Add("Subject", "1")

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		got, err := validation.ValidateBody[[]handlers.GetTicketResponse](response.Body)
		testutil.AssertValidResponse(t, err)

		testutil.AssertEqual(t, got, want)
	})

	t.Run("returns open tickets on GET to /tickets/open/", func(t *testing.T) {
		want := []handlers.GetTicketResponse{
			{
				ID: 1,
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
				ID: 2,
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
				ID: 3,
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
}
