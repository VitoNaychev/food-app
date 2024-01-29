package handlers

import (
	"net/http"

	"github.com/VitoNaychev/food-app/reqbuilder"
)

func NewGetLocationRequest(jwt string) *http.Request {
	request, _ := http.NewRequest(http.MethodGet, "/delivery/location/", nil)
	request.Header.Add("Token", jwt)

	return request
}

func NewUpdateLocationRequest(jwt string, lat float32, lon float32) *http.Request {
	updateLocationRequest := UpdateLocationRequest{
		Lat: lat,
		Lon: lon,
	}

	request := reqbuilder.NewRequestWithBody(http.MethodPost, "/delivery/location/", updateLocationRequest)
	request.Header.Add("Token", jwt)

	return request
}
