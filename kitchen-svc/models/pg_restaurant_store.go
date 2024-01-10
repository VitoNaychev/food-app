package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/jackc/pgx/v5"
)

type PgRestaurantStore struct {
	conn *pgx.Conn
}

func NewPgRestaurantStore(ctx context.Context, connString string) (*PgRestaurantStore, error) {
	conn, err := pgx.Connect(ctx, connString)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	pgRestaurantStore := PgRestaurantStore{conn}

	return &pgRestaurantStore, nil
}

func (p *PgRestaurantStore) DeleteRestaurant(id int) error {
	query := `DELETE FROM restaurants WHERE id = @id`
	args := pgx.NamedArgs{"id": id}

	_, err := p.conn.Exec(context.Background(), query, args)
	if err != nil {
		return storeerrors.FromPgxError(err)
	}

	return nil
}

func (p *PgRestaurantStore) CreateRestaurant(restaurant *Restaurant) error {
	createRestaurantQuery := `insert into restaurants(id) values (@id)`
	createRestaurantArgs := pgx.NamedArgs{
		"id": restaurant.ID,
	}

	tx, err := p.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), createRestaurantQuery, createRestaurantArgs)
	if err != nil {
		return storeerrors.FromPgxError(err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return storeerrors.FromPgxError(err)
	}

	return nil
}

func (p *PgRestaurantStore) GetRestaurantByID(id int) (Restaurant, error) {
	restaurantQuery := `select * from restaurants where id=@id`
	restaurantArgs := pgx.NamedArgs{
		"id": id,
	}

	row, _ := p.conn.Query(context.Background(), restaurantQuery, restaurantArgs)
	restaurant, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Restaurant])

	if err != nil {
		return Restaurant{}, storeerrors.FromPgxError(err)
	}

	return restaurant, nil
}
