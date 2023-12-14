package handlers

import "errors"

var (
	ErrExistingCourier    = errors.New("courier with this email already exists")
	ErrCourierNotFound    = errors.New("courier doesn't exists")
	ErrUnathorizedAction  = errors.New("courier does not have permission to perform this action")
	ErrInvalidCredentials = errors.New("invalid courier credentials")
)
