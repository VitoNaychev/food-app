package models

type Status int

const (
	CREATED     Status = 1 << 0
	ADDRESS_SET Status = 1 << 1
	HOURS_SET   Status = 1 << 2
	VALID       Status = CREATED | ADDRESS_SET | HOURS_SET
)
