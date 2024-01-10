package models

import (
	"context"
	"fmt"

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
