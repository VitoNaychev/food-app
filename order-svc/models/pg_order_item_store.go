package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/jackc/pgx/v5"
)

type PgOrderItemStore struct {
	conn *pgx.Conn
}

func NewPgOrderItemStore(ctx context.Context, connString string) (*PgOrderItemStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &PgOrderItemStore{conn}, nil
}

func (p *PgOrderItemStore) CreateOrderItem(orderItem *OrderItem) error {
	query := `insert into order_items(order_id, menu_item_id, quantity) 
	values (@order_id, @menu_item_id, @quantity) returning id`
	args := pgx.NamedArgs{
		"order_id":     orderItem.OrderID,
		"menu_item_id": orderItem.MenuItemID,
		"quantity":     orderItem.Quantity,
	}

	err := p.conn.QueryRow(context.Background(), query, args).Scan(&orderItem.ID)
	return storeerrors.FromPgxError(err)
}

func (p *PgOrderItemStore) GetOrderItemsByOrderID(orderID int) ([]OrderItem, error) {
	query := `select * from order_items where order_id=@order_id`
	args := pgx.NamedArgs{
		"order_id": orderID,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	orderItems, err := pgx.CollectRows(row, pgx.RowToStructByName[OrderItem])

	if err != nil {
		return []OrderItem{}, storeerrors.FromPgxError(err)
	}

	return orderItems, nil
}
