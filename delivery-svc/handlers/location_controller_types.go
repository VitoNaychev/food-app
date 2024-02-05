package handlers

import "github.com/VitoNaychev/food-app/delivery-svc/models"

type UpdateLocationRequest struct {
	Lat float32 `validate:"latitude,required"   json:"lat"`
	Lon float32 `validate:"longitude,required"  json:"lon"`
}

type GetLocationResponse struct {
	CourierID int     `validate:"min=1,required"   json:"courier_id"`
	Lat       float32 `validate:"latitude,required"   json:"lat"`
	Lon       float32 `validate:"longitude,required"  json:"lon"`
}

func LocationToGetLocationResponse(location models.Location) GetLocationResponse {
	getLocationResponse := GetLocationResponse{
		CourierID: location.CourierID,
		Lat:       location.Lat,
		Lon:       location.Lon,
	}
	return getLocationResponse
}
