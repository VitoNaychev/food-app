package orderevents

import "github.com/VitoNaychev/food-app/kitchen-svc/models"

var (
	shackTicket = models.Ticket{
		ID:           1,
		RestaurantID: 1,
		Total:        13.12,
		State:        models.OPEN,
		ReadyBy:      models.ZeroTime,
	}

	shackTicketItems = []models.TicketItem{
		{
			ID:         1,
			TicketID:   1,
			MenuItemID: 1,
			Quantity:   1,
		},
		{
			ID:         2,
			TicketID:   1,
			MenuItemID: 2,
			Quantity:   1,
		},
		{
			ID:         3,
			TicketID:   1,
			MenuItemID: 3,
			Quantity:   1,
		},
	}
)
