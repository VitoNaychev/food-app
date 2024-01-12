package stubs

import (
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/storeerrors"
)

type StubCourierStore struct {
	Couriers         []models.Courier
	CreatedCourier   models.Courier
	DeletedCourierID int
}

func (s *StubCourierStore) GetCourierByID(id int) (models.Courier, error) {
	for _, courier := range s.Couriers {
		if courier.ID == id {
			return courier, nil
		}
	}

	return models.Courier{}, storeerrors.ErrNotFound
}

func (s *StubCourierStore) CreateCourier(courier *models.Courier) error {
	s.CreatedCourier = *courier
	return nil
}

func (s *StubCourierStore) DeleteCourier(id int) error {
	s.DeletedCourierID = id
	return nil
}
