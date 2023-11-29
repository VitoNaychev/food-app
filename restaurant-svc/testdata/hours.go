package testdata

import (
	"time"

	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

var opening, _ = time.Parse("15:04", "10:00")
var closing, _ = time.Parse("15:04", "18:00")
var closed, _ = time.Parse("15:04", "00:00")

var (
	ShackHours = []models.Hours{
		{
			ID:           1,
			Day:          1,
			Opening:      opening,
			Closing:      closing,
			RestaurantID: 1,
		},
		{
			ID:           2,
			Day:          2,
			Opening:      opening,
			Closing:      closing,
			RestaurantID: 1,
		},
		{
			ID:           3,
			Day:          3,
			Opening:      opening,
			Closing:      closing,
			RestaurantID: 1,
		},
		{
			ID:           4,
			Day:          4,
			Opening:      opening,
			Closing:      closing,
			RestaurantID: 1,
		},
		{
			ID:           5,
			Day:          5,
			Opening:      opening,
			Closing:      closing,
			RestaurantID: 1,
		},
		{
			ID:           6,
			Day:          6,
			Opening:      opening,
			Closing:      closing,
			RestaurantID: 1,
		},
		{
			ID:           7,
			Day:          7,
			Opening:      closed,
			Closing:      closed,
			RestaurantID: 1,
		},
	}
	DominosHours = []models.Hours{
		{
			ID:           8,
			Day:          1,
			Opening:      closed,
			Closing:      closed,
			RestaurantID: 2,
		},
		{
			ID:           9,
			Day:          2,
			Opening:      opening,
			Closing:      closing,
			RestaurantID: 2,
		},
		{
			ID:           10,
			Day:          3,
			Opening:      opening,
			Closing:      closing,
			RestaurantID: 2,
		},
		{
			ID:           11,
			Day:          4,
			Opening:      opening,
			Closing:      closing,
			RestaurantID: 2,
		},
		{
			ID:           12,
			Day:          5,
			Opening:      opening,
			Closing:      closing,
			RestaurantID: 2,
		},
		{
			ID:           13,
			Day:          6,
			Opening:      opening,
			Closing:      closing,
			RestaurantID: 2,
		},
		{
			ID:           14,
			Day:          7,
			Opening:      opening,
			Closing:      closing,
			RestaurantID: 2,
		},
	}
)
