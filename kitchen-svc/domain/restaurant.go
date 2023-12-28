package domain

type Restaurant struct {
	ID int `db:"id"`
}

func NewRestaurant(id int) (Restaurant, error) {
	if id <= 0 {
		return Restaurant{}, ErrInvalidID
	}

	return Restaurant{ID: id}, nil
}
