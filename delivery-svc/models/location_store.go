package models

type LocationStore interface {
	CreateLocation(*Location) error
	GetLocationByCourierID(int) (Location, error)
	UpdateLocation(*Location) error
	DeleteLocation(int) error
}
