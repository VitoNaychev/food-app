package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/validation"
)

type LocationServer struct {
	secretKey []byte
	verifier  auth.Verifier

	locationStore models.LocationStore
	courierStore  models.CourierStore
}

func NewLocationServer(secretKey []byte, locationStore models.LocationStore, courierStore models.CourierStore) *LocationServer {
	locationServer := LocationServer{
		secretKey: secretKey,
		verifier:  NewCourierVerifier(courierStore),

		locationStore: locationStore,
		courierStore:  courierStore,
	}

	return &locationServer
}

func (l *LocationServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		auth.AuthenticationMW(l.updateLocation, l.verifier, l.secretKey)(w, r)
	case http.MethodGet:
		auth.AuthenticationMW(l.getLocation, l.verifier, l.secretKey)(w, r)
	}
}

func (l *LocationServer) updateLocation(w http.ResponseWriter, r *http.Request) {
	courierID, _ := strconv.Atoi(r.Header.Get("Subject"))

	updateLocationRequest, err := validation.ValidateBody[UpdateLocationRequest](r.Body)
	if err != nil {
		httperrors.HandleBadRequest(w, err)
		return
	}

	location := models.Location{
		CourierID: courierID,
		Lat:       updateLocationRequest.Lat,
		Lon:       updateLocationRequest.Lon,
	}
	err = l.locationStore.UpdateLocation(&location)
	if err != nil {
		httperrors.HandleBadRequest(w, err)
		return
	}

	json.NewEncoder(w).Encode(location)
}

func (l *LocationServer) getLocation(w http.ResponseWriter, r *http.Request) {
	courierID, _ := strconv.Atoi(r.Header.Get("Subject"))

	location, err := l.locationStore.GetLocationByCourierID(courierID)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
	}

	json.NewEncoder(w).Encode(location)
}
