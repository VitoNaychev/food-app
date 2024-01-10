package handlers

import "errors"

var (
	ErrUnathorizedAction         = errors.New("restaurant doesn't have permission to perform this action")
	ErrUnsuportedStateTransition = errors.New("current ticket state doesn't support this transition")
	ErrNonexistentState          = errors.New("such state doesn't exist")
)
