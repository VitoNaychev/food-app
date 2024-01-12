package handlers

import "errors"

var (
	ErrNoActiveDeliveries = errors.New("courier doesn't have active deliveries")
)
