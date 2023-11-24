package auth

import "errors"

var (
	ErrMissingToken       = errors.New("missing token")
	ErrInvalidCredentials = errors.New("invalid user credentials")
	ErrMissingSubject     = errors.New("token does not contain subject field")
	ErrNonIntegerSubject  = errors.New("token subject field is not an integer")
)
