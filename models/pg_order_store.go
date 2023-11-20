package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PgOrderStore struct {
	conn *pgx.Conn
}

func NewPgOrderStore(ctx context.Context, connString string) (*PgOrderStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &PgOrderStore{conn}, nil
}

func (p *PgOrderStore) GetOrdersByCustomerID(customerId int) ([]Order, error) {
	query := `select * from orders where customer_id=@customer_id`
	args := pgx.NamedArgs{
		"customer_id": customerId,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	orders, err := pgx.CollectRows(row, pgx.RowToStructByName[Order])

	if err != nil {
		return []Order{}, pgxErrorToStoreError(err)
	}

	return orders, nil
}

func (p *PgOrderStore) GetCurrentOrdersByCustomerID(customerId int) ([]Order, error) {
	query := `select * from orders where customer_id=@customer_id and status != @status`
	args := pgx.NamedArgs{
		"customer_id": customerId,
		"status":      COMPLETED,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	orders, err := pgx.CollectRows(row, pgx.RowToStructByName[Order])

	if err != nil {
		return []Order{}, pgxErrorToStoreError(err)
	}

	return orders, nil
}

func (p *PgOrderStore) CreateOrder(order *Order) error {
	query := `insert into orders(customer_id, restaurant_id, items, total, delivery_time, status, pickup_address, delivery_address) 
		values (@customer_id, @restaurant_id, @items, @total, @delivery_time, @status, @pickup_address, @delivery_address) returning id`
	args := pgx.NamedArgs{
		"customer_id":      order.CustomerID,
		"restaurant_id":    order.RestaurantID,
		"items":            order.Items,
		"total":            order.Total,
		"delivery_time":    order.DeliveryTime,
		"status":           order.Status,
		"pickup_address":   order.PickupAddress,
		"delivery_address": order.DeliveryAddress,
	}

	err := p.conn.QueryRow(context.Background(), query, args).Scan(&order.ID)
	return err

}
