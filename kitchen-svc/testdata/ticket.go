package testdata

import "github.com/VitoNaychev/food-app/kitchen-svc/models"

var (
	OpenShackTicket = models.Ticket{
		ID:           1,
		State:        models.CREATED,
		RestaurantID: 1,
	}

	InProgressShackTicket = models.Ticket{
		ID:           2,
		State:        models.PREPARING,
		RestaurantID: 1,
	}

	ReadyForPickupShackTicket = models.Ticket{
		ID:           3,
		State:        models.READY_FOR_PICKUP,
		RestaurantID: 1,
	}
)
