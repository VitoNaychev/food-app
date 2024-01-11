package models

type TicketItem struct {
	ID         int
	TicketID   int `db:"ticket_id"`
	MenuItemID int `db:"menu_item_id"`
	Quantity   int
}
