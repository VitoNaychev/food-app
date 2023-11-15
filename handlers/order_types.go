package handlers

type AuthStatus int

const (
	INVALID AuthStatus = iota
	NOT_FOUND
	OK
)

type AuthResponse struct {
	Status AuthStatus
	ID     int
}
