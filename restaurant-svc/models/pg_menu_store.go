package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/jackc/pgx/v5"
)

type PgMenuStore struct {
	conn *pgx.Conn
}

func NewPgMenuStore(ctx context.Context, connString string) (PgMenuStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return PgMenuStore{}, fmt.Errorf("unable to connect to database: %w", err)
	}

	pgMenuStore := PgMenuStore{conn}

	return pgMenuStore, nil
}

func (p *PgMenuStore) CreateMenuItem(menuItem *MenuItem) error {
	query := `insert into menu_items(name, price, details, restaurant_id) 
	values (@name, @price, @details, @restaurant_id) returning id`
	args := pgx.NamedArgs{
		"id":            menuItem.ID,
		"name":          menuItem.Name,
		"price":         menuItem.Price,
		"details":       menuItem.Details,
		"restaurant_id": menuItem.RestaurantID,
	}

	err := p.conn.QueryRow(context.Background(), query, args).Scan(&menuItem.ID)
	return err
}

func (p *PgMenuStore) GetMenuItemByID(id int) (MenuItem, error) {
	query := `select * from menu_items where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	menuItem, err := pgx.CollectOneRow(row, pgx.RowToStructByName[MenuItem])

	if err != nil {
		return MenuItem{}, storeerrors.FromPgxError(err)
	}

	return menuItem, nil
}

func (p *PgMenuStore) GetMenuByRestaurantID(restaurantID int) ([]MenuItem, error) {
	query := `select * from menu_items where restaurant_id=@restaurant_id`
	args := pgx.NamedArgs{
		"restaurant_id": restaurantID,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	menuItem, err := pgx.CollectRows(row, pgx.RowToStructByName[MenuItem])

	if err != nil {
		return []MenuItem{}, storeerrors.FromPgxError(err)
	}

	return menuItem, nil
}

func (p *PgMenuStore) DeleteMenuItem(id int) error {
	query := `delete from menu_items where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}

func (p *PgMenuStore) UpdateMenuItem(menuItem *MenuItem) error {
	query := `update menu_items set id=@id, name=@name, price=@price, 
	details=@details, restaurant_id=@restaurant_id where id=@id`
	args := pgx.NamedArgs{
		"id":            menuItem.ID,
		"name":          menuItem.Name,
		"price":         menuItem.Price,
		"details":       menuItem.Details,
		"restaurant_id": menuItem.RestaurantID,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}
