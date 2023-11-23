package handlers

import "errors"

var (
	ErrInvalidToken     = errors.New("token is invalid")
	ErrCustomerNotFound = errors.New("customer doesn't exist")
	ErrOrderNotFound    = errors.New("order doesn't exist")
)
