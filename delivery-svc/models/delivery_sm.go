package models

import "github.com/VitoNaychev/food-app/sm"

var deltas = []sm.Delta{
	{Current: sm.State(CREATED), Event: sm.Event(CANCEL_DELIVERY), Next: sm.State(CANCELED), Predicate: nil, Callback: nil},
	{Current: sm.State(CREATED), Event: sm.Event(DECLINE_DELIVERY), Next: sm.State(DECLINED), Predicate: nil, Callback: nil},
	{Current: sm.State(CREATED), Event: sm.Event(BEGIN_PREPARING_DELIVERY), Next: sm.State(IN_PROGRESS), Predicate: nil, Callback: nil},
	{Current: sm.State(IN_PROGRESS), Event: sm.Event(FINISH_PREPARING_DELIVERY), Next: sm.State(READY_FOR_PICKUP), Predicate: nil, Callback: nil},
	{Current: sm.State(READY_FOR_PICKUP), Event: sm.Event(PICKUP_DELIVERY), Next: sm.State(ON_ROUTE), Predicate: nil, Callback: nil},
	{Current: sm.State(ON_ROUTE), Event: sm.Event(COMPLETED), Next: sm.State(COMPLETED), Predicate: nil, Callback: nil},
}
