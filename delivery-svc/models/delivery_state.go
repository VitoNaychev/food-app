package models

import "errors"

var ErrNonexistentState = errors.New("such state doesn't exist")

type DeliveryState int

const (
	PENDING DeliveryState = iota
	CANCELED
	DECLINED
	IN_PROGRESS
	READY_FOR_PICKUP
	ON_ROUTE
	COMPLETED
)

func StateNameToStateValue(stateName string) (DeliveryState, error) {
	var stateValue DeliveryState

	switch stateName {
	case "pending":
		stateValue = PENDING
	case "canceled":
		stateValue = CANCELED
	case "in_progress":
		stateValue = IN_PROGRESS
	case "ready_for_pickup":
		stateValue = READY_FOR_PICKUP
	case "on_route":
		stateValue = ON_ROUTE
	case "completed":
		stateValue = COMPLETED
	default:
		return DeliveryState(-1), ErrNonexistentState
	}
	return stateValue, nil
}

func StateValueToStateName(stateValue DeliveryState) (string, error) {
	var stateName string

	switch stateValue {
	case PENDING:
		stateName = "pending"
	case CANCELED:
		stateName = "canceled"
	case IN_PROGRESS:
		stateName = "in_progress"
	case READY_FOR_PICKUP:
		stateName = "ready_for_pickup"
	case ON_ROUTE:
		stateName = "on_route"
	case COMPLETED:
		stateName = "completed"
	default:
		return "", ErrNonexistentState
	}

	return stateName, nil
}

type DeliveryEvent int

const (
	CANCEL_DELIVERY DeliveryEvent = iota
	DECLINE_DELIVERY
	BEGIN_PREPARING_DELIVERY
	FINISH_PREPARING_DELIVERY
	PICKUP_DELIVERY
	COMPLETE_DELIVERY
)
