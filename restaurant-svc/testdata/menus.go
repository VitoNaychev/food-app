package testdata

import "github.com/VitoNaychev/food-app/restaurant-svc/models"

var (
	ShackMenu = []models.MenuItem{
		{
			ID:           1,
			Name:         "XXL Duner",
			Price:        8.99,
			Details:      "The only thing you'll ever need",
			RestaurantID: 1,
		},
	}

	DominosMenu = []models.MenuItem{
		{
			ID:           1,
			Name:         "Burger Pizza",
			Price:        15.99,
			Details:      "Best pizza bruh",
			RestaurantID: 2,
		},
		{
			ID:           2,
			Name:         "Peperoni Pizza",
			Price:        13.99,
			Details:      "The OG pizza bruh",
			RestaurantID: 2,
		},
		{
			ID:           3,
			Name:         "Giros Pizza",
			Price:        14.99,
			Details:      "The new comer bruh",
			RestaurantID: 2,
		},
	}

	ForeignMenuItem = models.MenuItem{
		ID:           4,
		Name:         "Rizzoto",
		Price:        11.99,
		Details:      "Just a basic risoto",
		RestaurantID: 5,
	}
)
