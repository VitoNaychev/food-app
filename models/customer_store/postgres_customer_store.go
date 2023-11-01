package customer_store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type PostgresCustomerStore struct {
	conn *pgx.Conn
}

func NewPostgresCustomerStore(ctx context.Context, connString string) (PostgresCustomerStore, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return PostgresCustomerStore{}, fmt.Errorf("unable to connect to database: %w", err)
	}

	pgCustomerStore := PostgresCustomerStore{conn}
	return pgCustomerStore, nil
}

func (p *PostgresCustomerStore) DeleteCustomer(id int) error {
	panic("unimplemented")
}

func (p *PostgresCustomerStore) UpdateCustomer(customer Customer) error {
	panic("unimplemented")
}

func (p *PostgresCustomerStore) StoreCustomer(customer Customer) int {
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
	_ = p.conn.QueryRow(context.Background(), query, args).Scan(&customerId)

	return customerId
}

func (p *PostgresCustomerStore) GetCustomerByEmail(email string) (Customer, error) {
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

func (p *PostgresCustomerStore) GetCustomerById(id int) (Customer, error) {
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
