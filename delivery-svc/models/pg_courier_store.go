package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/jackc/pgx/v5"
)

type PgCourierStore struct {
	conn *pgx.Conn
}

func NewPgCourierStore(ctx context.Context, connString string) (*PgCourierStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &PgCourierStore{conn}, nil
}

func (p *PgCourierStore) GetCourierByID(id int) (Courier, error) {
	query := `select * from couriers where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	courier, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Courier])

	if err != nil {
		return Courier{}, storeerrors.FromPgxError(err)
	}

	return courier, nil
}

func (p *PgCourierStore) CreateCourier(courier *Courier) error {
	query := `insert into couriers(id, name) values (@id, @name)`
	args := pgx.NamedArgs{
		"id":   courier.ID,
		"name": courier.Name,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}

func (p *PgCourierStore) DeleteCourier(id int) error {
	query := `delete from couriers where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}
