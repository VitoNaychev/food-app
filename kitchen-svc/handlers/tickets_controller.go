package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/kitchen-svc/models"
	"github.com/VitoNaychev/food-app/validation"
)

type TicketServer struct {
	secretKey []byte

	ticketStore     models.TicketStore
	ticketItemStore models.TicketItemStore
	menuItemStore   models.MenuItemStore
	restaurantStore models.RestaurantStore

	verifier auth.Verifier
}

func NewTicketServer(secretKey []byte,
	ticketStore models.TicketStore,
	ticketItemStore models.TicketItemStore,
	menuItemStore models.MenuItemStore,
	restaurantStore models.RestaurantStore) *TicketServer {

	s := TicketServer{
		secretKey: secretKey,

		ticketStore:     ticketStore,
		ticketItemStore: ticketItemStore,
		menuItemStore:   menuItemStore,
		restaurantStore: restaurantStore,

		verifier: NewRestaurantVerifier(restaurantStore),
	}

	return &s
}

func (t *TicketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		auth.AuthenticationMW(t.getFilteredTickets, t.verifier, t.secretKey)(w, r)
	case http.MethodPost:
		auth.AuthenticationMW(t.stateTransitionHandler, t.verifier, t.secretKey)(w, r)
	}
}

func (t *TicketServer) stateTransitionHandler(w http.ResponseWriter, r *http.Request) {
	ticketRequest, _ := validation.ValidateBody[StateTransitionTicketRequest](r.Body)

	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))
	ticket, _ := t.ticketStore.GetTicketByID(ticketRequest.ID)

	if ticket.RestaurantID != restaurantID {
		httperrors.HandleUnauthorized(w, ErrUnathorizedAction)
		return
	}

	ticketSM := models.NewTicketSM(ticket.State)
	err := ticketSM.Exec(ticketRequest.Event)
	if err != nil {
		httperrors.HandleBadRequest(w, ErrUnsuportedStateTransition)
		return
	}

	_ = t.ticketStore.UpdateTicketState(ticketRequest.ID, ticketSM.Current())
}

func (t *TicketServer) getFilteredTickets(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	var tickets []models.Ticket
	stateName := r.URL.Query().Get("state")

	if stateName == "" {
		tickets, _ = t.ticketStore.GetTicketsByRestaurantID(restaurantID)
	} else {
		state, err := stateNameToStateValue(stateName)
		if err != nil {
			httperrors.HandleBadRequest(w, ErrNonexistentState)
			return
		}
		tickets, _ = t.ticketStore.GetTicketsByRestaurantIDWhereState(restaurantID, state)
	}

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

	json.NewEncoder(w).Encode(getTicketResponseArr)
}

func stateNameToStateValue(stateName string) (models.TicketState, error) {
	var stateValue models.TicketState

	switch stateName {
	case "open":
		stateValue = models.CREATED
	case "in_progress":
		stateValue = models.IN_PROGRESS
	case "ready_for_pickup":
		stateValue = models.READY_FOR_PICKUP
	case "completed":
		stateValue = models.COMPLETED
	case "declined":
		stateValue = models.DECLINED
	case "canceled":
		stateValue = models.CANCELED
	default:
		return models.TicketState(-1), ErrNonexistentState
	}

	return stateValue, nil
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
