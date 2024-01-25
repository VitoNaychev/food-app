package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/jackc/pgx/v5"
)

type PgTicketStore struct {
	conn *pgx.Conn
}

func NewPgTicketStore(ctx context.Context, connString string) (*PgTicketStore, error) {
	conn, err := pgx.Connect(ctx, connString)

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	pgTicketStore := PgTicketStore{conn}

	return &pgTicketStore, nil
}

func (p *PgTicketStore) CreateTicket(ticket *Ticket) error {
	query := `insert into tickets(id, restaurant_id, state, total, ready_by) 
	values (@id, @restaurant_id, @state, @total, @ready_by)`

	args := pgx.NamedArgs{
		"id":            ticket.ID,
		"restaurant_id": ticket.RestaurantID,
		"state":         ticket.State,
		"total":         ticket.Total,
		"ready_by":      ticket.ReadyBy,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}

func (p *PgTicketStore) GetTicketByID(id int) (Ticket, error) {
	query := `select * from tickets where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	tickets, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Ticket])

	if err != nil {
		return Ticket{}, storeerrors.FromPgxError(err)
	}

	return tickets, nil
}

func (p *PgTicketStore) GetTicketsByRestaurantID(restaurantID int) ([]Ticket, error) {
	query := `select * from tickets where restaurant_id=@restaurant_id`
	args := pgx.NamedArgs{
		"restaurant_id": restaurantID,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	tickets, err := pgx.CollectRows(row, pgx.RowToStructByName[Ticket])

	if err != nil {
		if err == pgx.ErrNoRows {
			return []Ticket{}, nil
		}
		return nil, err
	}

	return tickets, nil
}

func (p *PgTicketStore) GetTicketsByRestaurantIDWhereState(restaurantID int, state TicketState) ([]Ticket, error) {
	query := `select * from tickets where restaurant_id=@restaurant_id and state=@state`
	args := pgx.NamedArgs{
		"restaurant_id": restaurantID,
		"state":         state,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	tickets, err := pgx.CollectRows(row, pgx.RowToStructByName[Ticket])

	if err != nil {
		if err == pgx.ErrNoRows {
			return []Ticket{}, nil
		}
		return nil, err
	}

	return tickets, nil
}

func (p *PgTicketStore) UpdateTicket(ticket *Ticket) error {
	query := "update tickets set state=@state, ready_by=@ready_by, total=@total where id=@id "
	args := pgx.NamedArgs{
		"id":       ticket.ID,
		"state":    ticket.State,
		"total":    ticket.Total,
		"ready_by": ticket.ReadyBy,
	}

	_, err := p.conn.Exec(context.Background(), query, args)

	return storeerrors.FromPgxError(err)
}

func (p *PgTicketStore) UpdateTicketState(id int, state TicketState) error {
	query := "update tickets set state=@state where id=@id "
	args := pgx.NamedArgs{
		"id":    id,
		"state": state,
	}

	_, err := p.conn.Exec(context.Background(), query, args)

	return storeerrors.FromPgxError(err)
}
