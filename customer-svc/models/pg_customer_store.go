package models

import (
	"context"
	"fmt"

	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/jackc/pgx/v5"
)

type PgCustomerStore struct {
	conn *pgx.Conn
}

func NewPgCustomerStore(ctx context.Context, connString string) (PgCustomerStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return PgCustomerStore{}, fmt.Errorf("unable to connect to database: %w", err)
	}

	pgCustomerStore := PgCustomerStore{conn}
	return pgCustomerStore, nil
}

func (p *PgCustomerStore) GetCustomerByEmail(email string) (Customer, error) {
	query := `select * from customers where email=@email`
	args := pgx.NamedArgs{
		"email": email,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	customer, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Customer])

	if err != nil {
		return Customer{}, storeerrors.FromPgxError(err)
	}

	return customer, nil
}

func (p *PgCustomerStore) GetCustomerByID(id int) (Customer, error) {
	query := `select * from customers where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	customer, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Customer])

	if err != nil {
		return Customer{}, storeerrors.FromPgxError(err)
	}

	return customer, nil
}

func (p *PgCustomerStore) CreateCustomer(customer *Customer) error {
	query := `insert into customers(first_name, last_name, email, phone_number, password) 
		values (@firstName, @lastName, @email, @phone_number, @password) returning id`
	args := pgx.NamedArgs{
		"firstName":    customer.FirstName,
		"lastName":     customer.LastName,
		"email":        customer.Email,
		"phone_number": customer.PhoneNumber,
		"password":     customer.Password,
	}

	err := p.conn.QueryRow(context.Background(), query, args).Scan(&customer.Id)
	return storeerrors.FromPgxError(err)
}

func (p *PgCustomerStore) DeleteCustomer(id int) error {
	query := `delete from customers where id=@id`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}

func (p *PgCustomerStore) UpdateCustomer(customer *Customer) error {
	query := `update customers set first_name=@first_name, last_name=@last_name, 
		email=@email, phone_number=@phone_number, password=@password where id=@id`
	args := pgx.NamedArgs{
		"id":           customer.Id,
		"first_name":   customer.FirstName,
		"last_name":    customer.LastName,
		"email":        customer.Email,
		"phone_number": customer.PhoneNumber,
		"password":     customer.Password,
	}

	_, err := p.conn.Exec(context.Background(), query, args)
	return storeerrors.FromPgxError(err)
}
