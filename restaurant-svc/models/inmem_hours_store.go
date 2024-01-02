package models

import "github.com/VitoNaychev/food-app/storeerrors"

type InMemoryHoursStore struct {
	hours []Hours
}

func NewInMemoryHoursStore() *InMemoryHoursStore {
	return &InMemoryHoursStore{[]Hours{}}
}

func (i *InMemoryHoursStore) CreateHours(hour *Hours) error {
	hour.ID = len(i.hours) + 1
	i.hours = append(i.hours, *hour)
	return nil
}

func (i *InMemoryHoursStore) GetHoursByRestaurantID(restaurantID int) ([]Hours, error) {
	hours := []Hours{}
	for _, hour := range i.hours {
		if hour.RestaurantID == restaurantID {
			hours = append(hours, hour)
		}
	}
	return hours, nil
}

func (i *InMemoryHoursStore) UpdateHours(hour *Hours) error {
	for j, oldHour := range i.hours {
		if oldHour.ID == hour.ID {
			i.hours[j] = *hour
			return nil
		}
	}
	return storeerrors.ErrNotFound
}
