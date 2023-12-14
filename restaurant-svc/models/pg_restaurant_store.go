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

func NewPgRestaurantStore(ctx context.Context, connString string) (PgRestaurantStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return PgRestaurantStore{}, fmt.Errorf("unable to connect to database: %w", err)
	}

	pgRestaurantStore := PgRestaurantStore{conn}

	return pgRestaurantStore, nil
}

func (p *PgRestaurantStore) GetRestaurantByEmail(email string) (Restaurant, error) {
	query := `select * from restaurants where email=@email`
	args := pgx.NamedArgs{
		"email": email,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	restaurant, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Restaurant])

	if err != nil {
		return Restaurant{}, storeerrors.FromPgxError(err)
	}

	return restaurant, nil
}

func (p *PgRestaurantStore) GetRestaurantByID(id int) (Restaurant, error) {
	query := `select * from restaurants where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	restaurant, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Restaurant])

	if err != nil {
		return Restaurant{}, storeerrors.FromPgxError(err)
	}

	return restaurant, nil
}

func (p *PgRestaurantStore) CreateRestaurant(restaurant *Restaurant) error {
	query := `insert into restaurants(name, phone_number, email, password, IBAN, status) 
		values (@name, @phone_number, @email, @password, @iban, @status) returning id`
	args := pgx.NamedArgs{
		"name":         restaurant.Name,
		"phone_number": restaurant.PhoneNumber,
		"email":        restaurant.Email,
		"password":     restaurant.Password,
		"iban":         restaurant.IBAN,
		"status":       CREATED,
	}

	err := p.conn.QueryRow(context.Background(), query, args).Scan(&restaurant.ID)
	return storeerrors.FromPgxError(err)
}

func (p *PgRestaurantStore) DeleteRestaurant(id int) error {
	query := `delete from restaurants where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}

func (p *PgRestaurantStore) UpdateRestaurant(restaurant *Restaurant) error {
	query := `update restaurants set name=@name, phone_number=@phone_number, 
	email=@email, password=@password, IBAN=@iban, status=@status where id=@id`
	args := pgx.NamedArgs{
		"id":           restaurant.ID,
		"name":         restaurant.Name,
		"phone_number": restaurant.PhoneNumber,
		"email":        restaurant.Email,
		"password":     restaurant.Password,
		"iban":         restaurant.IBAN,
		"status":       restaurant.Status,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}
