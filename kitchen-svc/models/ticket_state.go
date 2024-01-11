package models

import "errors"

type TicketState int

var ErrNonexistentState = errors.New("such state doesn't exist")

const (
	CREATE_PENDING TicketState = iota
	REJECTED
	OPEN
	CANCELED
	DECLINED
	IN_PROGRESS
	READY_FOR_PICKUP
	COMPLETE_PENDING
	COMPLETED
)

func StateNameToStateValue(stateName string) (TicketState, error) {
	var stateValue TicketState

	switch stateName {
	case "open":
		stateValue = OPEN
	case "in_progress":
		stateValue = IN_PROGRESS
	case "ready_for_pickup":
		stateValue = READY_FOR_PICKUP
	case "completed":
		stateValue = COMPLETED
	case "declined":
		stateValue = DECLINED
	case "canceled":
		stateValue = CANCELED
	default:
		return TicketState(-1), ErrNonexistentState
	}
	return stateValue, nil
}

func StateValueToStateName(stateValue TicketState) (string, error) {
	var stateName string

	switch stateValue {
	case OPEN:
		stateName = "open"
	case IN_PROGRESS:
		stateName = "in_progress"
	case READY_FOR_PICKUP:
		stateName = "ready_for_pickup"
	case COMPLETED:
		stateName = "completed"
	case DECLINED:
		stateName = "declined"
	case CANCELED:
		stateName = "canceled"
	default:
		return "", ErrNonexistentState
	}

	return stateName, nil
}

type TicketEvent int

const (
	REJECT_TICKET TicketEvent = iota
	APPROVE_TICKET
	CANCEL_TICKET
	DECLINE_TICKET
	BEGIN_PREPARING
	FINISH_PREPARING
	COMPLETE_TICKET
	COMPLETE_REJECTED
	COMPLETE_CONFIRMED
)
