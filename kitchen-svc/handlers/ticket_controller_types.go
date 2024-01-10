package handlers

type StateTransitionTicketRequest struct {
	ID int
}

type GetTicketResponse struct {
	ID    int
	Items []GetTicketItemResponse
}

type GetTicketItemResponse struct {
	Quantity int
	Name     string
}
