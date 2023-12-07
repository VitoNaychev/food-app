package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/VitoNaychev/food-app/reqbuilder"
	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

func NewUpdateHoursRequest(jwt string, hours []models.Hours) *http.Request {
	updateHoursRequestArr := []HoursRequest{}
	for _, hour := range hours {
		updateHoursRequestArr = append(updateHoursRequestArr, HoursToHoursRequest(hour))
	}

	body := bytes.NewBuffer([]byte{})
	json.NewEncoder(body).Encode(updateHoursRequestArr)

	request, _ := http.NewRequest(http.MethodPut, "/restaurant/hours/", body)
	request.Header.Add("Token", jwt)

	return request
}

func NewCreateHoursRequest(jwt string, hours []models.Hours) *http.Request {
	createHoursRequestArr := []HoursRequest{}
	for _, hour := range hours {
		createHoursRequestArr = append(createHoursRequestArr, HoursToHoursRequest(hour))
	}

	request := reqbuilder.NewRequestWithBody[[]HoursRequest](
		http.MethodPost, "/restaurant/hours/", createHoursRequestArr)
	request.Header.Add("Token", jwt)

	return request
}

func NewGetHoursRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/restaurant/hours/", nil)
	request.Header.Add("Token", jwt)

	return request
}
