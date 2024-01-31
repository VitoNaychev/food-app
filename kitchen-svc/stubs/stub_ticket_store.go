package stubs

import (
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
)

type StubTicketStore struct {
	Tickets []models.Ticket

	SpyTicket models.Ticket
}

func (s *StubTicketStore) CreateTicket(ticket *models.Ticket) error {
	s.SpyTicket = *ticket

	return nil
}

func (s *StubTicketStore) GetTicketByID(id int) (models.Ticket, error) {
	if id != s.SpyTicket.ID {
		return models.Ticket{}, storeerrors.ErrNotFound
	}

	return s.SpyTicket, nil
}

func (s *StubTicketStore) UpdateTicket(ticket *models.Ticket) error {
	s.SpyTicket = *ticket

	return nil
}

func (s *StubTicketStore) UpdateTicketState(id int, state models.TicketState) error {
	if id != s.SpyTicket.ID {
		return storeerrors.ErrNotFound
	}

	s.SpyTicket.State = state
	return nil
}

func (s *StubTicketStore) GetTicketsByRestaurantID(restaurantID int) ([]models.Ticket, error) {
	restaruantTickets := []models.Ticket{}

	for _, ticket := range s.Tickets {
		if ticket.RestaurantID == restaurantID {
			restaruantTickets = append(restaruantTickets, ticket)
		}
	}

	return restaruantTickets, nil
}

func (s *StubTicketStore) GetTicketsByRestaurantIDWhereState(restaurantID int, state models.TicketState) ([]models.Ticket, error) {
	restaruantTickets := []models.Ticket{}

	for _, ticket := range s.Tickets {
		if ticket.RestaurantID == restaurantID &&
			ticket.State == state {
			restaruantTickets = append(restaruantTickets, ticket)
		}
	}

	return restaruantTickets, nil
}
