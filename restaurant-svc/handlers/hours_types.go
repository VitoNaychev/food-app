package handlers

import (
	"time"

	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type HoursRequest struct {
	Day     int    `validate:"min=1,max=7"           json:"day"`
	Opening string `validate:"required,workinghours" json:"opening"`
	Closing string `validate:"required,workinghours" json:"closing"`
}

func HoursToHoursRequest(hours models.Hours) HoursRequest {
	createHoursRequest := HoursRequest{
		Day:     hours.Day,
		Opening: hours.Opening.Format("15:04"),
		Closing: hours.Closing.Format("15:04"),
	}

	return createHoursRequest
}

func HoursRequestToHours(createHoursRequest HoursRequest, restaurantID int) models.Hours {
	opening, _ := time.Parse("15:04", createHoursRequest.Opening)
	closing, _ := time.Parse("15:04", createHoursRequest.Closing)
	hours := models.Hours{
		Day:          createHoursRequest.Day,
		Opening:      opening,
		Closing:      closing,
		RestaurantID: restaurantID,
	}

	return hours
}

type HoursResponse struct {
	ID           int    `validate:"min=1"                 json:"id"`
	Day          int    `validate:"min=1,max=7"           json:"day"`
	Opening      string `validate:"required,workinghours" json:"opening"`
	Closing      string `validate:"required,workinghours" json:"closing"`
	RestaurantID int    `validate:"min=1"                 json:"restaurant_id"`
}

func HoursToHoursResponse(hours models.Hours) HoursResponse {
	hoursResponse := HoursResponse{
		ID:           hours.ID,
		Day:          hours.Day,
		Opening:      hours.Opening.Format("15:04"),
		Closing:      hours.Closing.Format("15:04"),
		RestaurantID: hours.RestaurantID,
	}

	return hoursResponse
}

func HoursArrToHoursResponseArr(hoursArr []models.Hours) []HoursResponse {
	hoursResponseArr := []HoursResponse{}
	for _, hours := range hoursArr {
		hoursResponse := HoursToHoursResponse(hours)
		hoursResponseArr = append(hoursResponseArr, hoursResponse)
	}

	return hoursResponseArr
}
