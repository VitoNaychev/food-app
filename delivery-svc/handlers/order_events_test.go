package handlers_test

import (
	"testing"

	"github.com/VitoNaychev/food-app/delivery-svc/handlers"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/delivery-svc/stubs"
	"github.com/VitoNaychev/food-app/delivery-svc/testdata"

	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
	"github.com/VitoNaychev/food-app/testutil"
)

func TestOrderCreatedEventHandler(t *testing.T) {
	deliveryStore := &stubs.StubDeliveryStore{}
	addressStore := &stubs.StubAddressStore{CreatedAddresses: []models.Address{}}

	eventHandler := handlers.NewOrderEventHandler(deliveryStore, addressStore)

	t.Run("creates corresponding delivery", func(t *testing.T) {
		wantDelivery := testdata.VolenDelivery
		wantPickupAddress := testdata.VolenPickupAddress
		wantDeliveryAddress := testdata.VolenDeliveryAddress

		payload := testdata.PeterOrderCreatedEvent
		event := events.NewTypedEvent(svcevents.ORDER_CREATED_EVENT_ID, testdata.PeterOrderCreatedEvent.ID, payload)

		err := eventHandler.HandleOrderCreatedEvent(event)

		testutil.AssertNoErr(t, err)
		testutil.AssertEqual(t, deliveryStore.CreatedDelivery, wantDelivery)

		if len(addressStore.CreatedAddresses) != 2 {
			t.Fatalf("handler created %d addresses, want 2", len(addressStore.CreatedAddresses))
		}
		testutil.AssertEqual(t, addressStore.CreatedAddresses[0], wantPickupAddress)
		testutil.AssertEqual(t, addressStore.CreatedAddresses[1], wantDeliveryAddress)
	})
}
