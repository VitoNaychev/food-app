package models

import "github.com/VitoNaychev/food-app/storeerrors"

type InMemoryTicketStore struct {
	tickets []Ticket
}

func NewInMemoryTicketStore() *InMemoryTicketStore {
	return &InMemoryTicketStore{[]Ticket{}}
}

func (i *InMemoryTicketStore) CreateTicket(ticket *Ticket) error {
	i.tickets = append(i.tickets, *ticket)

	return nil
}

func (i *InMemoryTicketStore) GetTicketsByRestaurantID(restaurantID int) ([]Ticket, error) {
	var restaurantTickets []Ticket
	for _, ticket := range i.tickets {
		if ticket.RestaurantID == restaurantID {
			restaurantTickets = append(restaurantTickets, ticket)
		}
	}

	return restaurantTickets, nil
}

func (i *InMemoryTicketStore) GetTicketsByRestaurantIDWhereState(restaurantID int, state TicketState) ([]Ticket, error) {
	var filteredTickets []Ticket
	for _, ticket := range i.tickets {
		if ticket.RestaurantID == restaurantID && ticket.State == state {
			filteredTickets = append(filteredTickets, ticket)
		}
	}

	return filteredTickets, nil
}

func (i *InMemoryTicketStore) UpdateTicket(ticket *Ticket) error {
	for j, oldTicket := range i.tickets {
		if oldTicket.ID == ticket.ID {
			i.tickets[j] = *ticket
			return nil
		}
	}
	return storeerrors.ErrNotFound
}

func (i *InMemoryTicketStore) UpdateTicketState(ticketID int, newState TicketState) error {
	for j, ticket := range i.tickets {
		if ticket.ID == ticketID {
			i.tickets[j].State = newState
			return nil
		}
	}
	return storeerrors.ErrNotFound
}

func (i *InMemoryTicketStore) GetTicketByID(ticketID int) (Ticket, error) {
	for _, ticket := range i.tickets {
		if ticket.ID == ticketID {
			return ticket, nil
		}
	}
	return Ticket{}, storeerrors.ErrNotFound
}
