package models

type Restaurant struct {
	ID          int
	Name        string
	PhoneNumber string
	Email       string
	Password    string
	IBAN        string
	Status      Status
}
