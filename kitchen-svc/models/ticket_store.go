package models

type TicketStore interface {
	CreateTicket(*Ticket) error
	GetTicketsByRestaurantID(int) ([]Ticket, error)
	GetTicketsByRestaurantIDWhereState(int, TicketState) ([]Ticket, error)
	UpdateTicket(*Ticket) error
	UpdateTicketState(int, TicketState) error
	GetTicketByID(int) (Ticket, error)
}
