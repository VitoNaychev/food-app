package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PgOrderStore struct {
	conn *pgx.Conn
}

func NewPgOrderStore(ctx context.Context, connString string) (PgOrderStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return PgOrderStore{}, fmt.Errorf("unable to connect to database: %w", err)
	}

	PgOrderStore := PgOrderStore{conn}

	return PgOrderStore, nil
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
	query := `select * from orders where customer_id=@customer_id and status != 'COMPLETED's`
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
