package customer_store

type Customer struct {
	Id          int
	FirstName   string `db:"first_name"`
	LastName    string `db:"last_name"`
	PhoneNumber string `db:"phone_number"`
	Email       string
	Password    string
}
