package handlers

type DeleteMenuItemRequest struct {
	ID int `validate:"min=1"`
}

type UpdateMenuItemRequest struct {
	ID      int     `validate:"min=1"`
	Name    string  `validate:"min=2,max=20"`
	Price   float32 `validate:"required,price"`
	Details string  `validate:"max=1000"`
}

type CreateMenuItemRequest struct {
	Name    string  `validate:"min=2,max=20"`
	Price   float32 `validate:"required,price"`
	Details string  `validate:"max=1000"`
}

type GetMenuItemRequest struct {
	ID int `validate:"min=1"`
}
