package handlers

import "errors"

var (
	ErrMissingToken      = errors.New("token is missing")
	ErrInvalidToken      = errors.New("token is invalid")
	ErrCustomerNotFound  = errors.New("customer doesn't exist")
	ErrOrderNotFound     = errors.New("order doesn't exist")
	ErrUnathorizedAction = errors.New("customer does not have permission to perform this action")
)
