package stores

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/kitchen-svc/domain"
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

func (p *PgRestaurantStore) UpdateRestaurant(restaurant *domain.Restaurant) error {
	panic("unimplemented")
}

func (p *PgRestaurantStore) CreateRestaurant(restaurant *domain.Restaurant) error {
	createRestaurantQuery := `insert into restaurants(id) values (@id)`
	createRestaurantArgs := pgx.NamedArgs{
		"id": restaurant.ID,
	}

	// createMenuItemQuery := `insert into menu_items(id, restaurant_id, name, price)
	// 	values (@id, @restaurant_id, @name, @price)`

	tx, err := p.conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), createRestaurantQuery, createRestaurantArgs)
	if err != nil {
		return storeerrors.FromPgxError(err)
	}

	// for _, menuItem := range restaurant.MenuItems {
	// 	createMenuItemArgs := pgx.NamedArgs{
	// 		"id":            menuItem.ID,
	// 		"restaurant_id": menuItem.RestarautnID,
	// 		"name":          menuItem.Name,
	// 		"price":         menuItem.Price,
	// 	}

	// 	_, err = tx.Exec(context.Background(), createMenuItemQuery, createMenuItemArgs)
	// 	if err != nil {
	// 		return storeerrors.FromPgxError(err)
	// 	}
	// }

	err = tx.Commit(context.Background())
	if err != nil {
		return storeerrors.FromPgxError(err)
	}

	return nil
}

func (p *PgRestaurantStore) GetRestaurantByID(id int) (domain.Restaurant, error) {
	restaurantQuery := `select * from restaurants where id=@id`
	restaurantArgs := pgx.NamedArgs{
		"id": id,
	}

	// menuItemsQuery := `select * from menu_items where restaurant_id=@restaurant_id`
	// menuItemsArgs := pgx.NamedArgs{
	// 	"restaurant_id": id,
	// }

	row, _ := p.conn.Query(context.Background(), restaurantQuery, restaurantArgs)
	restaurant, err := pgx.CollectOneRow(row, pgx.RowToStructByName[domain.Restaurant])

	if err != nil {
		return domain.Restaurant{}, storeerrors.FromPgxError(err)
	}

	// rows, _ := p.conn.Query(context.Background(), menuItemsQuery, menuItemsArgs)
	// menuItems, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.MenuItem])

	// if err != nil {
	// 	return domain.Restaurant{}, storeerrors.FromPgxError(err)
	// }

	// restaurant.MenuItems = menuItems

	return restaurant, nil
}
