package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/kitchen-svc/models"
)

type TicketServer struct {
	ticketStore     models.TicketStore
	ticketItemStore models.TicketItemStore
	menuItemStore   models.MenuItemStore

	http.Handler
}

func NewTicketServer(ticketStore models.TicketStore, ticketItemStore models.TicketItemStore, menuItemStore models.MenuItemStore) *TicketServer {
	s := TicketServer{
		ticketStore:     ticketStore,
		ticketItemStore: ticketItemStore,
		menuItemStore:   menuItemStore,
	}

	router := http.NewServeMux()
	router.Handle("/tickets/open/", http.HandlerFunc(s.getOpenTickets))
	router.Handle("/tickets/in_progress/", http.HandlerFunc(s.getInProgressTickets))
	router.Handle("/tickets/ready_for_pickup/", http.HandlerFunc(s.getReadyForPickupTickets))

	s.Handler = router

	return &s
}

func (t *TicketServer) getOpenTickets(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	getTicketResponseArr, _ := t.getFilteredTickets(restaurantID, models.CREATED)

	json.NewEncoder(w).Encode(getTicketResponseArr)
}

func (t *TicketServer) getInProgressTickets(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	getTicketResponseArr, _ := t.getFilteredTickets(restaurantID, models.PREPARING)

	json.NewEncoder(w).Encode(getTicketResponseArr)
}

func (t *TicketServer) getReadyForPickupTickets(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	getTicketResponseArr, _ := t.getFilteredTickets(restaurantID, models.READY_FOR_PICKUP)

	json.NewEncoder(w).Encode(getTicketResponseArr)
}

func (t *TicketServer) getFilteredTickets(restaurantID int, state models.TicketState) ([]GetTicketResponse, error) {
	tickets, _ := t.ticketStore.GetTicketsByRestaurantIDWhereState(restaurantID, state)

	getTicketResponseArr := []GetTicketResponse{}

	for _, ticket := range tickets {
		ticketItems, _ := t.ticketItemStore.GetTicketItemsByTicketID(ticket.ID)

		getTicketItemResponseArr := []GetTicketItemResponse{}
		for _, ticketItem := range ticketItems {
			menuItem, _ := t.menuItemStore.GetMenuItemByID(ticketItem.MenuItemID)

			getTicketItemResponse := NewTicketItemResponse(ticketItem, menuItem)
			getTicketItemResponseArr = append(getTicketItemResponseArr, getTicketItemResponse)
		}

		getTicketResponse := NewTicketResponse(ticket, getTicketItemResponseArr)
		getTicketResponseArr = append(getTicketResponseArr, getTicketResponse)
	}

	return getTicketResponseArr, nil
}

func NewTicketItemResponse(ticketItem models.TicketItem, menuItem models.MenuItem) GetTicketItemResponse {
	getTicketItemResponse := GetTicketItemResponse{
		ID:       ticketItem.ID,
		Quantity: ticketItem.Quantity,
		Name:     menuItem.Name,
	}

	return getTicketItemResponse
}

func NewTicketResponse(ticket models.Ticket, getTicketItemResponseArr []GetTicketItemResponse) GetTicketResponse {
	getTicketResponse := GetTicketResponse{
		ID:    ticket.ID,
		Items: getTicketItemResponseArr,
	}

	return getTicketResponse
}
