package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/jackc/pgx/v5"
)

type PgMenuItemStore struct {
	conn *pgx.Conn
}

func NewPgMenuItemStore(ctx context.Context, connString string) (*PgMenuItemStore, error) {
	conn, err := pgx.Connect(ctx, connString)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	pgMenuItemStore := PgMenuItemStore{conn}

	return &pgMenuItemStore, nil
}

func (p *PgMenuItemStore) GetMenuItemByID(id int) (MenuItem, error) {
	query := `SELECT * FROM menu_items WHERE id = @id`
	args := pgx.NamedArgs{"id": id}

	row, _ := p.conn.Query(context.Background(), query, args)
	menuItem, err := pgx.CollectOneRow(row, pgx.RowToStructByName[MenuItem])

	if err != nil {
		return MenuItem{}, storeerrors.FromPgxError(err)
	}

	return menuItem, nil
}

func (p *PgMenuItemStore) CreateMenuItem(menuItem *MenuItem) error {
	query := `INSERT INTO menu_items (id, restaurant_id, name, price) 
	VALUES (@id, @restaurant_id, @name, @price) RETURNING id`
	args := pgx.NamedArgs{
		"id":            menuItem.ID,
		"restaurant_id": menuItem.RestaurantID,
		"name":          menuItem.Name,
		"price":         menuItem.Price,
	}

	err := p.conn.QueryRow(context.Background(), query, args).Scan(&menuItem.ID)

	if err != nil {
		return storeerrors.FromPgxError(err)
	}

	return nil
}

func (p *PgMenuItemStore) DeleteMenuItem(id int) error {
	query := `DELETE FROM menu_items WHERE id = @id`
	args := pgx.NamedArgs{"id": id}

	_, err := p.conn.Exec(context.Background(), query, args)
	if err != nil {
		return storeerrors.FromPgxError(err)
	}

	return nil
}

func (p *PgMenuItemStore) UpdateMenuItem(menuItem *MenuItem) error {
	query := `UPDATE menu_items SET restaurant_id = @restaurant_id, name = @name, price = @price WHERE id = @id`
	args := pgx.NamedArgs{
		"id":            menuItem.ID,
		"restaurant_id": menuItem.RestaurantID,
		"name":          menuItem.Name,
		"price":         menuItem.Price,
	}

	_, err := p.conn.Exec(context.Background(), query, args)

	if err != nil {
		return storeerrors.FromPgxError(err)
	}

	return nil
}

func (p *PgMenuItemStore) DeleteMenuItemWhereRestaurantID(restaurantID int) error {
	query := `DELETE FROM menu_items WHERE restaurant_id = @restaurant_id`
	args := pgx.NamedArgs{"restaurant_id": restaurantID}

	_, err := p.conn.Exec(context.Background(), query, args)

	if err != nil {
		return storeerrors.FromPgxError(err)
	}

	return nil
}
