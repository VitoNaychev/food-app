package handlers

import "github.com/VitoNaychev/food-app/restaurant-svc/models"

type UpdateAddressRequest struct {
	Lat          float64 `validate:"latitude,required"  json:"lat"`
	Lon          float64 `validate:"longitude,required" json:"lon"`
	AddressLine1 string  `validate:"required,max=40"    json:"address_line1"`
	AddressLine2 string  `validate:"max=40"             json:"address_line2"`
	City         string  `validate:"required,max=40"    json:"city"`
	Country      string  `validate:"required,max=20"    json:"country"`
}

func AddressToUpdateAddressRequest(address models.Address) UpdateAddressRequest {
	updateAddressRequest := UpdateAddressRequest{
		Lat:          address.Lat,
		Lon:          address.Lon,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		Country:      address.Country,
	}

	return updateAddressRequest
}

func UpdateAddressRequestToAddress(UpdateAddressRequest UpdateAddressRequest, id int, restauarntID int) models.Address {
	address := models.Address{
		ID:           id,
		RestaurantID: restauarntID,
		Lat:          UpdateAddressRequest.Lat,
		Lon:          UpdateAddressRequest.Lon,
		AddressLine1: UpdateAddressRequest.AddressLine1,
		AddressLine2: UpdateAddressRequest.AddressLine2,
		City:         UpdateAddressRequest.City,
		Country:      UpdateAddressRequest.Country,
	}

	return address
}

type DeleteAddressRequest struct {
	Id int `validate:"min=0" json:"id"`
}

type CreateAddressRequest struct {
	Lat          float64 `validate:"latitude,required"  json:"lat"`
	Lon          float64 `validate:"longitude,required" json:"lon"`
	AddressLine1 string  `validate:"required,max=40"    json:"address_line1"`
	AddressLine2 string  `validate:"max=40"             json:"address_line2"`
	City         string  `validate:"required,max=40"    json:"city"`
	Country      string  `validate:"required,max=20"    json:"country"`
}

func AddressToCreateAddressRequest(address models.Address) CreateAddressRequest {
	createAddressRequest := CreateAddressRequest{
		Lat:          address.Lat,
		Lon:          address.Lon,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		Country:      address.Country,
	}

	return createAddressRequest
}

func CreateAddressRequestToAddress(createAddressRequest CreateAddressRequest, restaurantID int) models.Address {
	address := models.Address{
		RestaurantID: restaurantID,
		Lat:          createAddressRequest.Lat,
		Lon:          createAddressRequest.Lon,
		AddressLine1: createAddressRequest.AddressLine1,
		AddressLine2: createAddressRequest.AddressLine2,
		City:         createAddressRequest.City,
		Country:      createAddressRequest.Country,
	}

	return address
}
