package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/jackc/pgx/v5"
)

type PgHoursStore struct {
	conn *pgx.Conn
}

func NewPgHoursStore(ctx context.Context, connString string) (PgHoursStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return PgHoursStore{}, fmt.Errorf("unable to connect to database: %w", err)
	}

	pgHoursStore := PgHoursStore{conn}

	return pgHoursStore, nil
}

func (p *PgHoursStore) CreateHours(hours *Hours) error {
	query := `insert into working_hours(day, opening, closing, restaurant_id) 
	values (@day, @opening, @closing, @restaurant_id) returning id`
	args := pgx.NamedArgs{
		"day":           hours.Day,
		"opening":       hours.Opening,
		"closing":       hours.Closing,
		"restaurant_id": hours.RestaurantID,
	}

	err := p.conn.QueryRow(context.Background(), query, args).Scan(&hours.ID)
	return err
}

func (p *PgHoursStore) GetHoursByRestaurantID(restaurantID int) ([]Hours, error) {
	query := `select * from working_hours where restaurant_id=@restaurant_id`
	args := pgx.NamedArgs{
		"restaurant_id": restaurantID,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	hours, err := pgx.CollectRows(row, pgx.RowToStructByName[Hours])

	if err != nil {
		return []Hours{}, storeerrors.FromPgxError(err)
	}

	return hours, nil
}

func (p *PgHoursStore) UpdateHours(hours *Hours) error {
	query := `update working_hours set day=@day, opening=@opening, closing=@closing  where id=@id`
	args := pgx.NamedArgs{
		"id":      hours.ID,
		"day":     hours.Day,
		"opening": hours.Opening,
		"closing": hours.Closing,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}
