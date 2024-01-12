package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/delivery-svc/models"
	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/storeerrors"
	"github.com/VitoNaychev/food-app/validation"
)

type DeliveryServer struct {
	deliveryStore models.DeliveryStore

	secretKey []byte
	verifier  auth.Verifier
}

func NewDeliveryServer(secretKey []byte, deliveryStore models.DeliveryStore, courierStore models.CourierStore) *DeliveryServer {
	deliveryServer := DeliveryServer{
		deliveryStore: deliveryStore,

		secretKey: secretKey,
		verifier:  NewCourierVerifier(courierStore),
	}

	return &deliveryServer
}

func (d *DeliveryServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		auth.AuthenticationMW(d.stateTransitionHandler, d.verifier, d.secretKey)(w, r)
	}
}

func (d *DeliveryServer) stateTransitionHandler(w http.ResponseWriter, r *http.Request) {
	courierID, _ := strconv.Atoi(r.Header.Get("Subject"))

	stateTransitionRequest, err := validation.ValidateBody[StateTransitionDeliveryRequest](r.Body)
	if err != nil {
		httperrors.HandleBadRequest(w, err)
		return
	}

	delivery, err := d.deliveryStore.GetActiveDeliveryByCourierID(courierID)
	if err != nil {
		if errors.Is(err, storeerrors.ErrNotFound) {
			httperrors.HandleBadRequest(w, ErrNoActiveDeliveries)
			return
		} else {
			httperrors.HandleInternalServerError(w, err)
			return
		}
	}

	deliverySM := models.NewDeliverySM(delivery.State)
	err = deliverySM.Exec(stateTransitionRequest.Event)
	if err != nil {
		httperrors.HandleBadRequest(w, err)
		return
	}

	delivery.State = deliverySM.Current()
	err = d.deliveryStore.UpdateDelivery(&delivery)
	if err != nil {
		httperrors.HandleInternalServerError(w, err)
		return
	}
}
