package models

type LocationStore interface {
	GeLocationByCourerID(int) (Location, error)
	UpdateLocation(*Location) error
}
