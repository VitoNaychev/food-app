package models

type LocationStore interface {
	GetLocationByCourerID(int) (Location, error)
	UpdateLocation(*Location) error
}
