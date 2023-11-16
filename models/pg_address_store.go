package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PgAddressStore struct {
	conn *pgx.Conn
}

func NewPgAddressStore(ctx context.Context, connString string) (PgAddressStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return PgAddressStore{}, fmt.Errorf("unable to connect to database: %w", err)
	}

	pgAddressStore := PgAddressStore{conn}

	return pgAddressStore, nil
}

func (p *PgAddressStore) GetAddressByID(id int) (Address, error) {
	query := `select * from addresses where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	address, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Address])

	if err != nil {
		return Address{}, pgxErrorToStoreError(err)
	}

	return address, nil
}
