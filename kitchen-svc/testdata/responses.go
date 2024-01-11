package testdata

import "github.com/VitoNaychev/food-app/kitchen-svc/handlers"

var (
	OpenShackTicketResponse = []handlers.GetTicketResponse{
		{
			ID:    1,
			State: "open",
			Total: 12.50,
			Items: []handlers.GetTicketItemResponse{
				{
					Quantity: 2,
					Name:     "Duner",
				},
			},
		},
	}

	InProgressShackTicketResponse = []handlers.GetTicketResponse{
		{
			ID:    2,
			State: "in_progress",
			Total: 25.00,
			Items: []handlers.GetTicketItemResponse{
				{
					Quantity: 4,
					Name:     "Duner",
				},
			},
		},
	}

	ReadyForPickupShackTicketResponse = []handlers.GetTicketResponse{
		{
			ID:    3,
			State: "ready_for_pickup",
			Total: 37.50,
			Items: []handlers.GetTicketItemResponse{
				{
					Quantity: 6,
					Name:     "Duner",
				},
			},
		},
	}

	CompletedShackTicketResponse = []handlers.GetTicketResponse{
		{
			ID:    4,
			State: "completed",
			Total: 6.25,
			Items: []handlers.GetTicketItemResponse{
				{
					Quantity: 1,
					Name:     "Duner",
				},
			},
		},
	}
)
