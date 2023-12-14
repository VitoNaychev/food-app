package handlers

import (
	"errors"

	"github.com/VitoNaychev/food-app/courier-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
)

type CourierVerifier struct {
	store models.CourierStore
}

func NewCourierVerifier(store models.CourierStore) *CourierVerifier {
	return &CourierVerifier{store}
}

func (c *CourierVerifier) DoesSubjectExist(id int) (bool, error) {
	_, err := c.store.GetCourierByID(id)
	if errors.Is(err, storeerrors.ErrNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
