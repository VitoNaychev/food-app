package handlers

import (
	"errors"

	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
)

type RestaurantVerifier struct {
	store models.RestaurantStore
}

func NewRestaurantVerifier(store models.RestaurantStore) *RestaurantVerifier {
	return &RestaurantVerifier{store}
}

func (c *RestaurantVerifier) DoesSubjectExist(id int) (bool, error) {
	_, err := c.store.GetRestaurantByID(id)
	if errors.Is(err, storeerrors.ErrNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
