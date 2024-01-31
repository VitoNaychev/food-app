package testdata

import "github.com/VitoNaychev/food-app/events/svcevents"

var (
	PeterOrderCreatedEventPickupAddress = svcevents.OrderCreatedEventAddress{
		ID:           1,
		Lat:          42.635934305,
		Lon:          23.380761684,
		AddressLine1: "ul. Filip Avramov 411, gk Mladost 4",
		AddressLine2: "",
		City:         "Sofia",
		Country:      "Bulgaria",
	}

	PeterOrderCreatedEventDeliveryAddress = svcevents.OrderCreatedEventAddress{
		ID:           2,
		Lat:          42.695111,
		Lon:          23.329184,
		AddressLine1: "Shipka Street 6",
		AddressLine2: "",
		City:         "Sofia",
		Country:      "Bulgaria",
	}

	PeterOrderCreatedEventItems = []svcevents.OrderCreatedEventItem{
		{
			ID:         1,
			MenuItemID: 1,
			Quantity:   2,
		},
	}

	PeterOrderCreatedEvent = svcevents.OrderCreatedEvent{
		ID:              1,
		RestaurantID:    1,
		Items:           PeterOrderCreatedEventItems,
		Total:           12.50,
		PickupAddress:   PeterOrderCreatedEventPickupAddress,
		DeliveryAddress: PeterOrderCreatedEventDeliveryAddress,
	}
)
