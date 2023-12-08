package handlers

import (
	"errors"
)

var (
	ErrExistingCustomer   = errors.New("customer with this email already exists")
	ErrCustomerNotFound   = errors.New("customer doesn't exists")
	ErrInvalidCredentials = errors.New("invalid user credentials")
	ErrMissingAddress     = errors.New("address doesn't exists")
	ErrUnathorizedAction  = errors.New("customer does not have permission to perform this action")
	ErrDatabaseError      = errors.New("operation encountered a database error")
)
