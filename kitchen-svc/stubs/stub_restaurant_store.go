package stubs

import "github.com/VitoNaychev/food-app/kitchen-svc/models"

type StubRestaurantStore struct {
	Restaurants         []models.Restaurant
	CreatedRestaurant   models.Restaurant
	DeletedRestaurantID int
}

func (s *StubRestaurantStore) DeleteRestaurant(id int) error {
	s.DeletedRestaurantID = id
	return nil
}

func (s *StubRestaurantStore) GetRestaurantByID(id int) (models.Restaurant, error) {
	for _, restaurant := range s.Restaurants {
		if restaurant.ID == id {
			return restaurant, nil
		}
	}

	return models.Restaurant{}, nil
}

func (s *StubRestaurantStore) CreateRestaurant(restaurant *models.Restaurant) error {
	s.CreatedRestaurant = *restaurant

	return nil
}
