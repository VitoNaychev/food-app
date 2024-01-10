package testdata

import "github.com/VitoNaychev/food-app/kitchen-svc/models"

var (
	OpenShackTicketItems = models.TicketItem{
		ID:         1,
		TicketID:   1,
		MenuItemID: 1,
		Quantity:   2,
	}
	InProgressShackTicketItems = models.TicketItem{
		ID:         2,
		TicketID:   2,
		MenuItemID: 1,
		Quantity:   4,
	}
	ReadyForPickupShackTicketItems = models.TicketItem{
		ID:         3,
		TicketID:   3,
		MenuItemID: 1,
		Quantity:   6,
	}
	PendingShackTicketItems = models.TicketItem{
		ID:         4,
		TicketID:   4,
		MenuItemID: 1,
		Quantity:   1,
	}
)
