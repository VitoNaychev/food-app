package models

type OrderItem struct {
	ID         int
	OrderID    int `db:"order_id"`
	MenuItemID int `db:"menu_item_id"`
	Quantity   int
}
