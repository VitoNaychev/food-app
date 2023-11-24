package msgtypes

type AuthStatus int

const (
	MISSING_TOKEN AuthStatus = iota
	INVALID
	NOT_FOUND
	OK
)

type AuthResponse struct {
	Status AuthStatus
	ID     int
}
