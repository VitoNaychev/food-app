package handlers

import (
	"time"

	"github.com/VitoNaychev/food-app/restaurant-svc/models"
)

type CreateHoursRequest struct {
	Day     int    `validate:"min=1,max=7"`
	Opening string `validate:"required,workinghours"`
	Closing string `validate:"required,workinghours"`
}

func HoursToCreateHoursRequest(hours models.Hours) CreateHoursRequest {
	createHoursRequest := CreateHoursRequest{
		Day:     hours.Day,
		Opening: hours.Opening.Format("15:04"),
		Closing: hours.Closing.Format("15:04"),
	}

	return createHoursRequest
}

func CreateHoursRequestToHours(createHoursRequest CreateHoursRequest, restaurantID int) models.Hours {
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
