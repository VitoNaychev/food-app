package models

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PgAddressStore struct {
	conn *pgx.Conn
}

func NewPgAddressStore(ctx context.Context, connString string) (*PgAddressStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &PgAddressStore{conn}, nil
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

func (p *PgAddressStore) CreateAddress(address *Address) error {
	query := `insert into addresses(lat, lon, address_line1, address_line2, city, country) 
	values (@lat, @lon, @address_line1, @address_line2, @city, @country) returning id`
	args := pgx.NamedArgs{
		"lat":           address.Lat,
		"lon":           address.Lon,
		"address_line1": address.AddressLine1,
		"address_line2": address.AddressLine2,
		"city":          address.City,
		"country":       address.Country,
	}

	err := p.conn.QueryRow(context.Background(), query, args).Scan(&address.ID)
	return err
}
