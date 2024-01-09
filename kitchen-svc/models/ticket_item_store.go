package models

type TicketItemStore interface {
	GetTicketItemsByTicketID(int) ([]TicketItem, error)
}
