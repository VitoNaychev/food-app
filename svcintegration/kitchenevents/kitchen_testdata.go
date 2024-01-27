package kitchenevents

import (
	"time"

	"github.com/VitoNaychev/food-app/kitchen-svc/models"
)

var (
	shackRestaurant = models.Restaurant{
		ID: 1,
	}

	shackMenuItem = models.MenuItem{
		ID:           1,
		RestaurantID: 1,
		Name:         "Duner",
		Price:        8.00,
	}

	shackTicket = models.Ticket{
		ID:           1,
		State:        models.OPEN,
		RestaurantID: 1,
		Total:        32.00,
		ReadyBy:      time.Time{},
	}

	shackTicketItem = models.TicketItem{
		ID:         1,
		TicketID:   1,
		MenuItemID: 1,
		Quantity:   4,
	}
)
