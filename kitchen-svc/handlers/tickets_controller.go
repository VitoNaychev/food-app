package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/validation"
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
	router.Handle("/tickets/accept/", http.HandlerFunc(s.acceptTicket))
	router.Handle("/tickets/decline/", http.HandlerFunc(s.declineTicket))
	router.Handle("/tickets/prepared/", http.HandlerFunc(s.preparedTicket))

	router.Handle("/tickets/open/", http.HandlerFunc(s.getOpenTickets))
	router.Handle("/tickets/in_progress/", http.HandlerFunc(s.getInProgressTickets))
	router.Handle("/tickets/ready_for_pickup/", http.HandlerFunc(s.getReadyForPickupTickets))
	router.Handle("/tickets/completed/", http.HandlerFunc(s.getCompletedTickets))

	s.Handler = router

	return &s
}

func (t *TicketServer) preparedTicket(w http.ResponseWriter, r *http.Request) {
	t.stateTransitionHandler(w, r, models.FINISH_PREPARING)
}

func (t *TicketServer) declineTicket(w http.ResponseWriter, r *http.Request) {
	t.stateTransitionHandler(w, r, models.DECLINE_TICKET)
}

func (t *TicketServer) acceptTicket(w http.ResponseWriter, r *http.Request) {
	t.stateTransitionHandler(w, r, models.START_PREPARING)
}

func (t *TicketServer) stateTransitionHandler(w http.ResponseWriter, r *http.Request, event models.TicketEvent) {
	ticketRequest, _ := validation.ValidateBody[StateTransitionTicketRequest](r.Body)

	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))
	ticket, _ := t.ticketStore.GetTicketByID(ticketRequest.ID)

	if ticket.RestaurantID != restaurantID {
		httperrors.HandleUnauthorized(w, ErrUnathorizedAction)
		return
	}

	ticketSM := models.NewTicketSM(ticket.State)
	err := ticketSM.Exec(event)
	if err != nil {
		httperrors.HandleBadRequest(w, ErrUnsuportedStateTransition)
		return
	}

	_ = t.ticketStore.UpdateTicketState(ticketRequest.ID, ticketSM.Current())
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

func (t *TicketServer) getCompletedTickets(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	getTicketResponseArr, _ := t.getFilteredTickets(restaurantID, models.COMPLETED)

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
		Quantity: ticketItem.Quantity,
		Name:     menuItem.Name,
	}

	return getTicketItemResponse
}

func NewTicketResponse(ticket models.Ticket, getTicketItemResponseArr []GetTicketItemResponse) GetTicketResponse {
	getTicketResponse := GetTicketResponse{
		ID:    ticket.ID,
		Total: ticket.Total,
		Items: getTicketItemResponseArr,
	}

	return getTicketResponse
}
