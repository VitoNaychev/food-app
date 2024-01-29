package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/jackc/pgx/v5"
)

type PgDeliveryStore struct {
	conn *pgx.Conn
}

func NewPgDeliveryStore(ctx context.Context, connString string) (*PgDeliveryStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &PgDeliveryStore{conn}, nil
}

func (p *PgDeliveryStore) GetDeliveryByID(id int) (Delivery, error) {
	query := `select * from deliveries where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	delivery, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Delivery])

	if err != nil {
		return Delivery{}, storeerrors.FromPgxError(err)
	}

	return delivery, nil
}

func (p *PgDeliveryStore) GetActiveDeliveryByCourierID(courierID int) (Delivery, error) {
	query := `select * from deliveries where courier_id=@courier_id and 
		(state != @pending or state != @canceled or state != @completed)`
	args := pgx.NamedArgs{
		"courier_id": courierID,
		"pending":    PENDING,
		"canceled":   CANCELED,
		"completed":  COMPLETED,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	delivery, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Delivery])

	if err != nil {
		return Delivery{}, storeerrors.FromPgxError(err)
	}

	return delivery, nil
}

func (p *PgDeliveryStore) UpdateDelivery(delivery *Delivery) error {
	query := `update deliveries set courier_id=@courier_id, pickup_address_id=@pickup_address_id, 
		delivery_address_id=@delivery_address_id, ready_by=@ready_by, state=@state where id=@id`
	args := pgx.NamedArgs{
		"id":                  delivery.ID,
		"courier_id":          delivery.CourierID,
		"pickup_address_id":   delivery.PickupAddressID,
		"delivery_address_id": delivery.DeliveryAddressID,
		"ready_by":            delivery.ReadyBy,
		"state":               delivery.State,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}
