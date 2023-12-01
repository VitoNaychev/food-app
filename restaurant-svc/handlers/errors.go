package handlers

import (
	"errors"
)

var (
	ErrMissingAddress     = errors.New("address doesn't exists")
	ErrExistingRestaurant = errors.New("restauarnt with this email already exists")
	ErrRestaurantNotFound = errors.New("restaurant doesn't exists")
	ErrHoursAlreadySet    = errors.New("hours for restaurant are already set")
	ErrAddressAlreadySet  = errors.New("address for restaurant is already set")
	ErrIncompleteWeek     = errors.New("working hours are not set for every day of the week")
	ErrDuplicateDays      = errors.New("trying to set duplicate working hours for the same day")
)
