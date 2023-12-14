package handlers

import (
	"errors"

	"github.com/VitoNaychev/food-app/customer-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
)

type CustomerVerifier struct {
	store models.CustomerStore
}

func NewCustomerVerifier(store models.CustomerStore) *CustomerVerifier {
	return &CustomerVerifier{store}
}

func (c *CustomerVerifier) DoesSubjectExist(id int) (bool, error) {
	_, err := c.store.GetCustomerByID(id)
	if errors.Is(err, storeerrors.ErrNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
