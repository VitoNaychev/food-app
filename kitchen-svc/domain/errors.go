package domain

import "errors"

var (
	ErrInvalidID = errors.New("trying to create an object with invalid ID")
)
