package models

type Restaurant struct {
	ID          int
	Name        string
	PhoneNumber string `db:"phone_number"`
	Email       string
	Password    string
	IBAN        string
	Status      Status
}
