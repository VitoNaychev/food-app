package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/events"
	"github.com/VitoNaychev/food-app/events/svcevents"
	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/VitoNaychev/food-app/validation"
)

func (s *CourierServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginCourierRequest, err := validation.ValidateBody[LoginCourierRequest](r.Body)
	if err != nil {
		httperrors.HandleBadRequest(w, err)
		return
	}

	courier, err := s.store.GetCourierByEmail(loginCourierRequest.Email)
	if err != nil {
		if errors.Is(err, storeerrors.ErrNotFound) {
			httperrors.HandleUnauthorized(w, ErrInvalidCredentials)
			return
		} else {
			httperrors.HandleInternalServerError(w, err)
			return
		}
	}

	if courier.Password != loginCourierRequest.Password {
		httperrors.HandleUnauthorized(w, ErrInvalidCredentials)
		return
	}

	jwtToken, _ := auth.GenerateJWT(s.secretKey, s.expiresAt, courier.ID)
	jwtResponse := JWTResponse{Token: jwtToken}

	json.NewEncoder(w).Encode(jwtResponse)
}

func (s *CourierServer) deleteCourier(w http.ResponseWriter, r *http.Request) {
	courierID, _ := strconv.Atoi(r.Header.Get("Subject"))

	err := s.store.DeleteCourier(courierID)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
	}

	payload := svcevents.CourierDeletedEvent{ID: courierID}
	event := events.NewEvent(svcevents.COURIER_CREATED_EVENT_ID, courierID, payload)

	err = s.publisher.Publish(svcevents.COURIER_EVENTS_TOPIC, event)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
	}
}

func (s *CourierServer) updateCourier(w http.ResponseWriter, r *http.Request) {
	courierID, _ := strconv.Atoi(r.Header.Get("Subject"))

	updateCourierRequest, err := validation.ValidateBody[UpdateCourierRequest](r.Body)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	newCourier := UpdateCourierRequestToCourier(updateCourierRequest, courierID)
	newCourier.ID = courierID

	err = s.store.UpdateCourier(&newCourier)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	updateCourierResponse := CourierToCourierResponse(newCourier)
	json.NewEncoder(w).Encode(updateCourierResponse)
}

func (s *CourierServer) getCourier(w http.ResponseWriter, r *http.Request) {
	courierID, _ := strconv.Atoi(r.Header.Get("Subject"))

	courier, err := s.store.GetCourierByID(courierID)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
	}

	getCourierResponse := CourierToCourierResponse(courier)
	json.NewEncoder(w).Encode(getCourierResponse)
}

func (s *CourierServer) createCourier(w http.ResponseWriter, r *http.Request) {
	createCourierRequest, err := validation.ValidateBody[CreateCourierRequest](r.Body)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	courier := CreateCourierRequestToCourier(createCourierRequest)
	if _, err = s.store.GetCourierByEmail(courier.Email); !errors.Is(err, storeerrors.ErrNotFound) {
		httperrors.WriteJSONError(w, http.StatusBadRequest, ErrExistingCourier)
		return
	}

	err = s.store.CreateCourier(&courier)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, err)
		return
	}

	jwtToken, _ := auth.GenerateJWT(s.secretKey, s.expiresAt, courier.ID)

	response := CreateCourierResponse{
		JWT:     JWTResponse{jwtToken},
		Courier: CourierToCourierResponse(courier),
	}
	json.NewEncoder(w).Encode(response)

	payload := svcevents.CourierCreatedEvent{
		ID:   courier.ID,
		Name: courier.FirstName,
	}
	event := events.NewEvent(svcevents.COURIER_CREATED_EVENT_ID, courier.ID, payload)

	err = s.publisher.Publish(svcevents.COURIER_EVENTS_TOPIC, event)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
	}
}
