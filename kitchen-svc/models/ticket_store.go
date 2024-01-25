package models

type TicketStore interface {
	GetTicketsByRestaurantID(int) ([]Ticket, error)
	GetTicketsByRestaurantIDWhereState(int, TicketState) ([]Ticket, error)
	UpdateTicket(*Ticket) error
	UpdateTicketState(int, TicketState) error
	GetTicketByID(int) (Ticket, error)
}
