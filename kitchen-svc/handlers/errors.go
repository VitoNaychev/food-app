package handlers

import "errors"

var (
	ErrUnathorizedAction         = errors.New("restaurant doesn't have permission to perform this action")
	ErrUnsuportedStateTransition = errors.New("current ticket state doesn't support this transition")
	ErrInvalidTimeFormat         = errors.New("time string has invalid format")
	ErrInvalidTime               = errors.New("time string must represent a moment in the future")
)
