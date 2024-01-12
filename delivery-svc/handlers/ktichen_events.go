package handlers

import (
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
)

type KitchenEventHandler struct {
	deliveryStore models.DeliveryStore
}

func NewKitchenEventHandler(deliveryStore models.DeliveryStore) *KitchenEventHandler {
	kitchenEventHandler := KitchenEventHandler{
		deliveryStore: deliveryStore,
	}

	return &kitchenEventHandler
}

func (k *KitchenEventHandler) HandleTicketBeginPreparingEvent(event events.Event[svcevents.TicketBeginPreparingEvent]) error {
	delivery, err := k.deliveryStore.GetDeliveryByID(event.Payload.ID)
	if err != nil {
		return err
	}

	err = k.applyEventToDelivery(&delivery, models.BEGIN_PREPARING_DELIVERY)
	if err != nil {
		return err
	}

	delivery.ReadyBy = event.Payload.ReadyBy
	err = k.deliveryStore.UpdateDelivery(&delivery)
	if err != nil {
		return err
	}

	return nil
}

func (k *KitchenEventHandler) HandleTicketCancelEvent(event events.Event[svcevents.TicketCancelEvent]) error {
	err := k.applyEventAndUpdateDelivery(event.Payload.ID, models.CANCEL_DELIVERY)

	return err
}

func (k *KitchenEventHandler) HandleTicketFinishPreparingEvent(event events.Event[svcevents.TicketFinishPreparingEvent]) error {
	err := k.applyEventAndUpdateDelivery(event.Payload.ID, models.FINISH_PREPARING_DELIVERY)

	return err
}

func (k *KitchenEventHandler) applyEventAndUpdateDelivery(deliveryId int, event models.DeliveryEvent) error {
	delivery, err := k.deliveryStore.GetDeliveryByID(deliveryId)
	if err != nil {
		return err
	}

	err = k.applyEventToDelivery(&delivery, event)
	if err != nil {
		return err
	}

	err = k.deliveryStore.UpdateDelivery(&delivery)
	if err != nil {
		return err
	}

	return nil
}

func (k *KitchenEventHandler) applyEventToDelivery(delivery *models.Delivery, event models.DeliveryEvent) error {
	deliverySM := models.NewDeliverySM(delivery.State)
	err := deliverySM.Exec(event)
	if err != nil {
		return err
	}

	delivery.State = deliverySM.Current()
	return nil
}
