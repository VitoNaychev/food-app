package models

import (
	"context"
	"fmt"

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

func (p *PgCustomerStore) DeleteCustomer(id int) error {
	panic("unimplemented")
}

func (p *PgCustomerStore) UpdateCustomer(customer *Customer) error {
	panic("unimplemented")
}

func (p *PgCustomerStore) CreateCustomer(customer *Customer) error {
	query := `insert into customers(first_name, last_name, email, phone_number, password) 
		values (@firstName, @lastName, @email, @phoneNumber, @password) returning id`
	args := pgx.NamedArgs{
		"firstName":   customer.FirstName,
		"lastName":    customer.LastName,
		"email":       customer.Email,
		"phoneNumber": customer.PhoneNumber,
		"password":    customer.Password,
	}

	var customerId int
	err := p.conn.QueryRow(context.Background(), query, args).Scan(&customerId)

	return err
}

func (p *PgCustomerStore) GetCustomerByEmail(email string) (Customer, error) {
	query := `select * from customers where email=@email`
	args := pgx.NamedArgs{
		"email": email,
	}

	row, _ := p.conn.Query(context.Background(), query, args)
	customer, err := pgx.CollectOneRow(row, pgx.RowToStructByName[Customer])

	if err != nil {
		return Customer{}, err
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
		return Customer{}, err
	}

	return customer, nil
}
