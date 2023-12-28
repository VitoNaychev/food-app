package services

import (
	"github.com/VitoNaychev/food-app/kitchen-svc/domain"
	"github.com/VitoNaychev/food-app/kitchen-svc/stores"
)

type KitchenService struct {
	restaurantStore stores.RestaurantStore
}

func NewKitchenService(restaurantStore stores.RestaurantStore) *KitchenService {
	service := KitchenService{
		restaurantStore: restaurantStore,
	}

	return &service
}

func (k *KitchenService) CreateRestaurant(id int) error {
	restaurant, err := domain.NewRestaurant(id)
	if err != nil {
		return err
	}

	err = k.restaurantStore.CreateRestaurant(&restaurant)
	if err != nil {
		return err
	}

	return nil
}

func (k *KitchenService) CreateMenuItem(id int, restaurantID int, name string, price float32) error {
	panic("unimplemented")
}
