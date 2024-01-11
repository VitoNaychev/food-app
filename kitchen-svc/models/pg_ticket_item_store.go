package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/jackc/pgx/v5"
)

type PgTicketItemStore struct {
	conn *pgx.Conn
}

func NewPgTicketItemStore(ctx context.Context, connString string) (*PgTicketItemStore, error) {
	conn, err := pgx.Connect(ctx, connString)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	pgTicketItemStore := PgTicketItemStore{conn}

	return &pgTicketItemStore, nil
}

func (p *PgTicketItemStore) CreateTicketItem(ticketItem *TicketItem) error {
	query := `insert into ticket_items(id, ticket_id, menu_item_id, quantity) 
	values (@id, @ticket_id, @menu_item_id, @quantity)`
	args := pgx.NamedArgs{
		"id":           ticketItem.ID,
		"ticket_id":    ticketItem.TicketID,
		"menu_item_id": ticketItem.MenuItemID,
		"quantity":     ticketItem.Quantity,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}

func (p *PgTicketItemStore) GetTicketItemsByTicketID(ticketID int) ([]TicketItem, error) {
	query := `select * from ticket_items where ticket_id=@ticket_id`
	args := pgx.NamedArgs{
		"ticket_id": ticketID,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	ticketItems, err := pgx.CollectRows(row, pgx.RowToStructByName[TicketItem])

	if err != nil {
		if err == pgx.ErrNoRows {
			return []TicketItem{}, nil
		}
		return nil, err
	}

	return ticketItems, nil
}
