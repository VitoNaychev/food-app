package models

import "github.com/VitoNaychev/food-app/sm"

var deliveryDeltas = []sm.Delta{
	{Current: sm.State(PENDING), Event: sm.Event(CANCEL_DELIVERY), Next: sm.State(CANCELED), Predicate: nil, Callback: nil},
	{Current: sm.State(PENDING), Event: sm.Event(DECLINE_DELIVERY), Next: sm.State(DECLINED), Predicate: nil, Callback: nil},
	{Current: sm.State(PENDING), Event: sm.Event(BEGIN_PREPARING_DELIVERY), Next: sm.State(IN_PROGRESS), Predicate: nil, Callback: nil},
	{Current: sm.State(IN_PROGRESS), Event: sm.Event(FINISH_PREPARING_DELIVERY), Next: sm.State(READY_FOR_PICKUP), Predicate: nil, Callback: nil},
	{Current: sm.State(READY_FOR_PICKUP), Event: sm.Event(PICKUP_DELIVERY), Next: sm.State(ON_ROUTE), Predicate: nil, Callback: nil},
	{Current: sm.State(ON_ROUTE), Event: sm.Event(COMPLETE_DELIVERY), Next: sm.State(COMPLETED), Predicate: nil, Callback: nil},
}

type DeliverySM struct {
	sm sm.SM
}

func NewDeliverySM(initial DeliveryState) DeliverySM {
	sm := sm.New(sm.State(initial), deliveryDeltas, nil)
	return DeliverySM{sm}
}

func (d *DeliverySM) Exec(event DeliveryEvent) error {
	err := d.sm.Exec(sm.Event(event))
	return err
}

func (d *DeliverySM) Current() DeliveryState {
	return DeliveryState(d.sm.Current)
}
