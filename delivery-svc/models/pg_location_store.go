package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/jackc/pgx/v5"
)

type PgLocationStore struct {
	conn *pgx.Conn
}

func NewPgLocationStore(ctx context.Context, connString string) (*PgLocationStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	return &PgLocationStore{conn}, nil
}

func (p *PgLocationStore) CreateLocation(location *Location) error {
	query := `insert into locations(courier_id, lat, lon) values (@courier_id, @lat, @lon)`
	args := pgx.NamedArgs{
		"courier_id": location.CourierID,
		"lat":        location.Lat,
		"lon":        location.Lon,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}

func (p *PgLocationStore) GetLocationByCourierID(courierID int) (Location, error) {
	query := `select * from locations where courier_id=@courier_id`
	args := pgx.NamedArgs{
		"courier_id": courierID,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	location, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Location])

	if err != nil {
		return Location{}, storeerrors.FromPgxError(err)
	}

	return location, nil
}

func (p *PgLocationStore) UpdateLocation(location *Location) error {
	query := `update locations set lat=@lat, lon=@lon where courier_id=@courier_id`
	args := pgx.NamedArgs{
		"courier_id": location.CourierID,
		"lat":        location.Lat,
		"lon":        location.Lon,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}

func (p *PgLocationStore) DeleteLocation(courierID int) error {
	query := `delete from locations where courier_id=@courier_id`
	args := pgx.NamedArgs{
		"courier_id": courierID,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}
