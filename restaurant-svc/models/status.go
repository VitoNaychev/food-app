package models

type Status int

const (
	CREATION_PENDING Status = iota
	ADDRESS_PENDING
	HOURS_PENDING
	VALID
)
