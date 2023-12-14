package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/jackc/pgx/v5"
)

type PgCourierStore struct {
	conn *pgx.Conn
}

func NewPgCourierStore(ctx context.Context, connString string) (PgCourierStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return PgCourierStore{}, fmt.Errorf("unable to connect to database: %w", err)
	}

	pgCourierStore := PgCourierStore{conn}

	return pgCourierStore, nil
}

func (p *PgCourierStore) GetCourierByEmail(email string) (Courier, error) {
	query := `select * from couriers where email=@email`
	args := pgx.NamedArgs{
		"email": email,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	courier, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Courier])

	if err != nil {
		return Courier{}, storeerrors.FromPgxError(err)
	}

	return courier, nil
}

func (p *PgCourierStore) GetCourierByID(id int) (Courier, error) {
	query := `select * from couriers where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	courier, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Courier])

	if err != nil {
		return Courier{}, storeerrors.FromPgxError(err)
	}

	return courier, nil
}

func (p *PgCourierStore) CreateCourier(courier *Courier) error {
	query := `insert into couriers(first_name, last_name, phone_number, email, password, IBAN) 
		values (@first_name, @last_name, @phone_number, @email, @password, @iban) returning id`
	args := pgx.NamedArgs{
		"first_name":   courier.FirstName,
		"last_name":    courier.LastName,
		"phone_number": courier.PhoneNumber,
		"email":        courier.Email,
		"password":     courier.Password,
		"iban":         courier.IBAN,
	}

	err := p.conn.QueryRow(context.Background(), query, args).Scan(&courier.ID)
	return storeerrors.FromPgxError(err)
}

func (p *PgCourierStore) DeleteCourier(id int) error {
	query := `delete from couriers where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}

func (p *PgCourierStore) UpdateCourier(courier *Courier) error {
	query := `update couriers set first_name=@first_name, last_name=@last_name, phone_number=@phone_number, 
	email=@email, password=@password, IBAN=@iban where id=@id`
	args := pgx.NamedArgs{
		"id":           courier.ID,
		"first_name":   courier.FirstName,
		"last_name":    courier.LastName,
		"phone_number": courier.PhoneNumber,
		"email":        courier.Email,
		"password":     courier.Password,
		"iban":         courier.IBAN,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}
