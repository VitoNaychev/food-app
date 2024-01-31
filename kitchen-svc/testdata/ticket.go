package testdata

import (
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
)

var (
	OpenShackTicket = models.Ticket{
		ID:           1,
		Total:        12.50,
		State:        models.OPEN,
		RestaurantID: 1,
		ReadyBy:      models.ZeroTime,
	}

	InProgressShackTicket = models.Ticket{
		ID:           2,
		Total:        25.00,
		State:        models.IN_PROGRESS,
		RestaurantID: 1,
		ReadyBy:      ReadyByTime,
	}

	ReadyForPickupShackTicket = models.Ticket{
		ID:           3,
		Total:        37.50,
		State:        models.READY_FOR_PICKUP,
		RestaurantID: 1,
		ReadyBy:      ReadyByTime,
	}

	CompletedShackTicket = models.Ticket{
		ID:           4,
		Total:        6.25,
		State:        models.COMPLETED,
		RestaurantID: 1,
		ReadyBy:      ReadyByTime,
	}

	ForeginRestaurantTicket = models.Ticket{
		ID:           5,
		Total:        10.00,
		State:        models.OPEN,
		RestaurantID: 5,
		ReadyBy:      models.ZeroTime,
	}
)
