package models

type DeliveryStatus int

const (
	CREATED DeliveryStatus = iota
	CANCELED
	DECLINED
	IN_PROGRESS
	READY_FOR_PICKUP
	ON_ROUTE
	COMPLETED
)

type DeliveryEvent int

const (
	CANCEL_DELIVERY DeliveryEvent = iota
	DECLINE_DELIVERY
	BEGIN_PREPARING_DELIVERY
	FINISH_PREPARING_DELIVERY
	PICKUP_DELIVERY
	COMPLETE_DELIVERY
)
