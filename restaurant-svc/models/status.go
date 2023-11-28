package models

type Status int

const (
	CREATION_PENDING Status = 0
	ADDRESS_SET      Status = 1 << 0
	HOURS_SET        Status = 1 << 1
	VALID            Status = ADDRESS_SET | HOURS_SET
)
