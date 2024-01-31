package stubs

import "github.com/VitoNaychev/food-app/kitchen-svc/models"

type StubTicketItemStore struct {
	TicketItems    []models.TicketItem
	SpyTicketItems []models.TicketItem
}

func (s *StubTicketItemStore) CreateTicketItem(ticketItem *models.TicketItem) error {
	s.SpyTicketItems = append(s.SpyTicketItems, *ticketItem)
	return nil
}

func (s *StubTicketItemStore) GetTicketItemsByTicketID(ticketID int) ([]models.TicketItem, error) {
	ticketItems := []models.TicketItem{}

	for _, ticketItem := range s.TicketItems {
		if ticketItem.TicketID == ticketID {
			ticketItems = append(ticketItems, ticketItem)
		}
	}

	return ticketItems, nil
}
