package models

type InMemoryTicketItemStore struct {
	ticketItems []TicketItem
}

func NewInMemoryTicketItemStore() *InMemoryTicketItemStore {
	return &InMemoryTicketItemStore{[]TicketItem{}}
}

func (i *InMemoryTicketItemStore) CreateTicketItem(ticketItem *TicketItem) error {
	i.ticketItems = append(i.ticketItems, *ticketItem)

	return nil
}

func (i *InMemoryTicketItemStore) GetTicketItemsByTicketID(ticketID int) ([]TicketItem, error) {
	var ticketItems []TicketItem

	for _, item := range i.ticketItems {
		if item.TicketID == ticketID {
			ticketItems = append(ticketItems, item)
		}
	}

	return ticketItems, nil
}
