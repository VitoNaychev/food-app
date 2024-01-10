package models

import "github.com/VitoNaychev/food-app/sm"

var ticketDeltas = []sm.Delta{
	{Current: sm.State(CREATE_PENDING), Event: sm.Event(APPROVE_TICKET), Next: sm.State(CREATED), Predicate: nil, Callback: nil},
	{Current: sm.State(CREATE_PENDING), Event: sm.Event(REJECT_TICKET), Next: sm.State(REJECTED), Predicate: nil, Callback: nil},
	{Current: sm.State(CREATED), Event: sm.Event(CANCEL_TICKET), Next: sm.State(CANCELED), Predicate: nil, Callback: nil},
	{Current: sm.State(CREATED), Event: sm.Event(DECLINE_TICKET), Next: sm.State(DECLINED), Predicate: nil, Callback: nil},
	{Current: sm.State(CREATED), Event: sm.Event(BEGIN_PREPARING), Next: sm.State(IN_PROGRESS), Predicate: nil, Callback: nil},
	{Current: sm.State(IN_PROGRESS), Event: sm.Event(FINISH_PREPARING), Next: sm.State(READY_FOR_PICKUP), Predicate: nil, Callback: nil},
	{Current: sm.State(READY_FOR_PICKUP), Event: sm.Event(COMPLETE_TICKET), Next: sm.State(COMPLETE_PENDING), Predicate: nil, Callback: nil},
	{Current: sm.State(COMPLETE_PENDING), Event: sm.Event(COMPLETE_REJECTED), Next: sm.State(READY_FOR_PICKUP), Predicate: nil, Callback: nil},
	{Current: sm.State(COMPLETE_PENDING), Event: sm.Event(COMPLETE_CONFIRMED), Next: sm.State(COMPLETED), Predicate: nil, Callback: nil},
}

type TicketSM struct {
	sm sm.SM
}

func NewTicketSM(initial TicketState) TicketSM {
	sm := sm.New(sm.State(initial), ticketDeltas, nil)
	return TicketSM{sm}
}

func (t *TicketSM) Exec(event TicketEvent) error {
	err := t.sm.Exec(sm.Event(event))
	return err
}

func (t *TicketSM) Current() TicketState {
	return TicketState(t.sm.Current)
}
