package handlers

import (
	"errors"

	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type RestaurantVerifier struct {
	store models.RestaurantStore
}

func NewRestaurantVerifier(store models.RestaurantStore) *RestaurantVerifier {
	return &RestaurantVerifier{store}
}

func (c *RestaurantVerifier) DoesSubjectExist(id int) (bool, error) {
	_, err := c.store.GetRestaurantByID(id)
	if errors.Is(err, models.ErrNotFound) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
