package handlers

import (
	"errors"
)

var (
	ErrMissingAddress     = errors.New("address doesn't exists")
	ErrExistingRestaurant = errors.New("restauarnt with this email already exists")
	ErrRestaurantNotFound = errors.New("restaurant doesn't exists")
	ErrAddressAlreadySet  = errors.New("address for restaurant is already set")
)
