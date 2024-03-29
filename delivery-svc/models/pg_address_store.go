package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
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
		return Address{}, storeerrors.FromPgxError(err)
	}

	return address, nil
}

func (p *PgAddressStore) CreateAddress(address *Address) error {
	query := `insert into addresses(id, lat, lon, address_line1, address_line2, city, country) 
	values (@id, @lat, @lon, @address_line1, @address_line2, @city, @country)`
	args := pgx.NamedArgs{
		"id":            address.ID,
		"lat":           address.Lat,
		"lon":           address.Lon,
		"address_line1": address.AddressLine1,
		"address_line2": address.AddressLine2,
		"city":          address.City,
		"country":       address.Country,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return err
}
