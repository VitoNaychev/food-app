package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/VitoNaychev/food-app/courier-svc/models"
	"github.com/VitoNaychev/food-app/reqbuilder"
)

func NewLoginCourierRequest(courier models.Courier) *http.Request {
	requestBody := LoginCourierRequest{courier.Email, courier.Password}
	request := reqbuilder.NewRequestWithBody[LoginCourierRequest](
		http.MethodPost, "/courier/login/", requestBody)

	return request
}

func NewDeleteCourierRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodDelete, "/courier/", nil)
	request.Header.Add("Token", jwt)

	return request
}

func NewUpdateCourierRequest(jwt string, courier models.Courier) *http.Request {
	updateCourierRequest := CourierToUpdateCourierRequest(courier)
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(updateCourierRequest)

	request, _ := http.NewRequest(http.MethodPut, "/courier/", body)
	request.Header.Add("Token", jwt)

	return request
}

func NewCreateCourierRequest(courier models.Courier) *http.Request {
	createCourierRequest := CourierToCreateCourierRequest(courier)
	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(createCourierRequest)

	request, _ := http.NewRequest(http.MethodPost, "/courier/", body)
	return request
}

func NewGetCourierRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/courier/", nil)
	request.Header.Add("Token", jwt)

	return request
}
