package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
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

	publisher events.EventPublisher

	verifier auth.Verifier
}

func NewTicketServer(secretKey []byte,
	ticketStore models.TicketStore,
	ticketItemStore models.TicketItemStore,
	menuItemStore models.MenuItemStore,
	restaurantStore models.RestaurantStore,
	publisher events.EventPublisher) *TicketServer {

	s := TicketServer{
		secretKey: secretKey,

		ticketStore:     ticketStore,
		ticketItemStore: ticketItemStore,
		menuItemStore:   menuItemStore,
		restaurantStore: restaurantStore,

		publisher: publisher,

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

	ticket, err := t.ticketStore.GetTicketByID(ticketRequest.ID)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
		return
	}

	if ticket.RestaurantID != restaurantID {
		httperrors.HandleUnauthorized(w, ErrUnathorizedAction)
		return
	}

	ticketSM := models.NewTicketSM(ticket.State)
	err = ticketSM.Exec(ticketRequest.Event)
	if err != nil {
		httperrors.HandleBadRequest(w, ErrUnsuportedStateTransition)
		return
	}

	ticket.State = ticketSM.Current()

	if ticketRequest.Event == models.BEGIN_PREPARING {
		readyBy, err := ParseTimeAndSetDate(ticketRequest.ReadyBy)
		if err != nil {
			httperrors.HandleBadRequest(w, err)
		}

		ticket.ReadyBy = readyBy
	}

	err = t.ticketStore.UpdateTicket(&ticket)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
		return
	}

	ticket.State = ticketSM.Current()
	stateTransitionResponse := NewStateTransitionResponse(ticket)

	json.NewEncoder(w).Encode(stateTransitionResponse)

	err = t.sendTicketStateTransitionEvent(ticketRequest.Event, ticket)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
	}
}

func (t *TicketServer) sendTicketStateTransitionEvent(event models.TicketEvent, ticket models.Ticket) error {
	switch event {
	case models.BEGIN_PREPARING:
		payload := svcevents.TicketBeginPreparingEvent{
			ID:      ticket.ID,
			ReadyBy: ticket.ReadyBy,
		}
		event := events.NewEvent(svcevents.TICKET_BEGIN_PREPARING_EVENT_ID, ticket.ID, payload)

		err := t.publisher.Publish(svcevents.KITCHEN_EVENTS_TOPIC, event)
		return err
	case models.FINISH_PREPARING:
		payload := svcevents.TicketFinishPreparingEvent{
			ID: ticket.ID,
		}
		event := events.NewEvent(svcevents.TICKET_FINISH_PREPARING_EVENT_ID, ticket.ID, payload)

		err := t.publisher.Publish(svcevents.KITCHEN_EVENTS_TOPIC, event)
		return err
	default:
		// no event needs to be published
		return nil
	}
}

func ParseTimeAndSetDate(readyByStr string) (time.Time, error) {
	readyByTime, err := time.Parse("15:04", readyByStr)
	if err != nil {
		return time.Time{}, ErrInvalidTimeFormat
	}

	now := time.Now()
	readyByDateTime := time.Date(now.Year(), now.Month(), now.Day(), readyByTime.Hour(), readyByTime.Minute(), 0, 0, now.Location())

	if !readyByDateTime.After(now) {
		return time.Time{}, ErrInvalidTime
	}

	return readyByDateTime, nil
}

func (t *TicketServer) getFilteredTickets(w http.ResponseWriter, r *http.Request) {
	restaurantID, _ := strconv.Atoi(r.Header.Get("Subject"))

	tickets, err := t.getTicketsForRestaurantAndQueryParams(restaurantID, r.URL.Query())
	if err != nil {
		if errors.Is(err, models.ErrNonexistentState) {
			httperrors.HandleBadRequest(w, err)
		} else {
			httperrors.HandleInternalServerError(w, err)
		}
	}

	getTicketResponseArr, err := t.newGetTicketResponseArrForTickets(tickets)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
	}

	json.NewEncoder(w).Encode(getTicketResponseArr)
}

func (t *TicketServer) getTicketsForRestaurantAndQueryParams(restaurantID int, params url.Values) ([]models.Ticket, error) {
	var tickets []models.Ticket

	stateName := params.Get("state")
	if stateName == "" {
		var err error

		tickets, err = t.ticketStore.GetTicketsByRestaurantID(restaurantID)
		if err != nil {
			return nil, err
		}
	} else {
		var state models.TicketState
		state, err := models.StateNameToStateValue(stateName)
		if err != nil {
			return nil, err
		}

		tickets, err = t.ticketStore.GetTicketsByRestaurantIDWhereState(restaurantID, state)
		if err != nil {
			return nil, err
		}
	}

	return tickets, nil
}

func (t *TicketServer) newGetTicketResponseArrForTickets(tickets []models.Ticket) ([]GetTicketResponse, error) {
	getTicketResponseArr := []GetTicketResponse{}

	for _, ticket := range tickets {
		ticketItems, err := t.ticketItemStore.GetTicketItemsByTicketID(ticket.ID)
		if err != nil {
			return nil, err
		}

		getTicketItemResponseArr := []GetTicketItemResponse{}
		for _, ticketItem := range ticketItems {
			menuItem, err := t.menuItemStore.GetMenuItemByID(ticketItem.MenuItemID)
			if err != nil {
				return nil, err
			}

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
	stateName, _ := models.StateValueToStateName(ticket.State)

	getTicketResponse := GetTicketResponse{
		ID:      ticket.ID,
		Total:   ticket.Total,
		State:   stateName,
		Items:   getTicketItemResponseArr,
		ReadyBy: ticket.ReadyBy,
	}

	return getTicketResponse
}
