package stubs

import (
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
)

type StubLocationStore struct {
	Locations                []models.Location
	UpdatedLocation          models.Location
	CreatedLocation          models.Location
	DeletedLocationCourierID int
}

func (s *StubLocationStore) DeleteLocation(courierID int) error {
	s.DeletedLocationCourierID = courierID

	return nil
}

func (s *StubLocationStore) CreateLocation(location *models.Location) error {
	s.CreatedLocation = *location

	return nil
}

func (s *StubLocationStore) GetLocationByCourierID(courierID int) (models.Location, error) {
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
