package stubs

import (
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
)

type StubLocationStore struct {
	Locations       []models.Location
	UpdatedLocation models.Location
}

func (s *StubLocationStore) GeLocationByCourerID(courierID int) (models.Location, error) {
	for _, location := range s.Locations {
		if location.CourierID == courierID {
			return location, nil
		}
	}

	return models.Location{}, storeerrors.ErrNotFound
}

func (s *StubLocationStore) UpdateLocation(location *models.Location) error {
	s.UpdatedLocation = *location

	return nil
}
