package models

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

type StoreError struct {
	message string
}

func (d *StoreError) Error() string {
	return d.message
}

func NewStoreError(message string) *StoreError {
	return &StoreError{message}
}

var (
	ErrNotFound = &StoreError{"didn't find object in database"}
)

func pgxErrorToStoreError(err error) error {
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}

		return NewStoreError(err.Error())
	}
	return nil
}
