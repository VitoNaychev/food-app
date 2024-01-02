package models

import "github.com/VitoNaychev/food-app/storeerrors"

type InMemoryRestaurantStore struct {
	restaurants []Restaurant
}

func NewInMemoryRestaurantStore() *InMemoryRestaurantStore {
	return &InMemoryRestaurantStore{[]Restaurant{}}
}

func (i *InMemoryRestaurantStore) CreateRestaurant(restaurant *Restaurant) error {
	restaurant.ID = len(i.restaurants) + 1
	i.restaurants = append(i.restaurants, *restaurant)

	return nil
}

func (i *InMemoryRestaurantStore) DeleteRestaurant(id int) error {
	for j, restaurant := range i.restaurants {
		if restaurant.ID == id {
			i.restaurants = append(i.restaurants[:j], i.restaurants[j+1:]...)
			return nil
		}
	}

	return storeerrors.ErrNotFound
}

func (i *InMemoryRestaurantStore) GetRestaurantByEmail(email string) (Restaurant, error) {
	for _, restaurant := range i.restaurants {
		if restaurant.Email == email {
			return restaurant, nil
		}
	}

	return Restaurant{}, storeerrors.ErrNotFound
}

func (i *InMemoryRestaurantStore) GetRestaurantByID(id int) (Restaurant, error) {
	for _, restaurant := range i.restaurants {
		if restaurant.ID == id {
			return restaurant, nil
		}
	}

	return Restaurant{}, storeerrors.ErrNotFound
}

func (i *InMemoryRestaurantStore) UpdateRestaurant(restaurant *Restaurant) error {
	for j, oldRestaurant := range i.restaurants {
		if oldRestaurant.ID == restaurant.ID {
			i.restaurants[j] = *restaurant
			return nil
		}
	}

	return storeerrors.ErrNotFound
}
