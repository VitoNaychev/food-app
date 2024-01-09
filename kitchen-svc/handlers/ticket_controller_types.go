package handlers

type GetTicketResponse struct {
	ID    int
	Items []GetTicketItemResponse
}

type GetTicketItemResponse struct {
	ID       int
	Quantity int
	Name     string
}
