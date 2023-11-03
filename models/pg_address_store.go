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

func (p *PgAddressStore) CreateAddress(address *Address) error {
	query := `insert into addresses(customer_id, lat, lon, address_line1, address_line2, city, country) 
	values (@customer_id, @lat, @lon, @address_line1, @address_line2, @city, @country) returning id`
	args := pgx.NamedArgs{
		"customer_id":   address.CustomerId,
		"lat":           address.Lat,
		"lon":           address.Lon,
		"address_line1": address.AddressLine1,
		"address_line2": address.AddressLine2,
		"city":          address.City,
		"country":       address.Country,
	}

	err := p.conn.QueryRow(context.Background(), query, args).Scan(&address.Id)
	return err
}

func (p *PgAddressStore) GetAddressByID(id int) (Address, error) {
	query := `select * from addresses where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	address, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Address])

	if err != nil {
		return Address{}, err
	}

	return address, nil
}
func (p *PgAddressStore) GetAddressesByCustomerID(customerID int) ([]Address, error) {
	query := `select * from addresses where customer_id=@customer_id`
	args := pgx.NamedArgs{
		"customer_id": customerID,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	address, err := pgx.CollectRows(row, pgx.RowToStructByName[Address])

	if err != nil {
		return []Address{}, err
	}

	return address, nil
}

func (p *PgAddressStore) DeleteAddress(id int) error {
	query := `delete from addresses where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return err
}

func (p *PgAddressStore) UpdateAddress(address *Address) error {
	query := `update addresses set lat=@lat, lon=@lon, address_line1=@address_line1,
	address_line2=@address_line2, city=@city, country=@country where id=@id`
	args := pgx.NamedArgs{
		"id":            address.Id,
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
