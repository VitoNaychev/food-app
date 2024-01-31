package models

type TicketItemStore interface {
	CreateTicketItem(*TicketItem) error
	GetTicketItemsByTicketID(int) ([]TicketItem, error)
}
