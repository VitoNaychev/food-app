package handlers

import (
	"errors"
)

var (
	ErrHoursNotSet        = errors.New("working hours are not set")
	ErrAddressNotSet      = errors.New("address is not set")
	ErrExistingRestaurant = errors.New("restauarnt with this email already exists")
	ErrRestaurantNotFound = errors.New("restaurant doesn't exists")
	ErrHoursAlreadySet    = errors.New("hours for restaurant are already set")
	ErrAddressAlreadySet  = errors.New("address for restaurant is already set")
	ErrIncompleteWeek     = errors.New("working hours are not set for every day of the week")
	ErrDuplicateDays      = errors.New("trying to set duplicate working hours for the same day")
	ErrMissingMenuItem    = errors.New("menu item doesn't exist")
	ErrUnathorizedAction  = errors.New("restaurant does not have permission to perform this action")
	ErrInvalidRestaurant  = errors.New("restaurant has hours or location not set")
)
