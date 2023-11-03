package address

import "github.com/VitoNaychev/bt-customer-svc/models"

type UpdateAddressRequest struct {
	Id           int     `validate:"min=0"`
	Lat          float64 `validate:"latitude,required"`
	Lon          float64 `validate:"longitude,required"`
	AddressLine1 string  `validate:"required,max=40"`
	AddressLine2 string  `validate:"max=40"`
	City         string  `validate:"required,max=40"`
	Country      string  `validate:"required,max=20"`
}

func AddressToUpdateAddressRequest(address models.Address) UpdateAddressRequest {
	updateAddressRequest := UpdateAddressRequest{
		Id:           address.Id,
		Lat:          address.Lat,
		Lon:          address.Lon,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		Country:      address.Country,
	}

	return updateAddressRequest
}

func UpdateAddressRequestToAddress(UpdateAddressRequest UpdateAddressRequest, customerId int) models.Address {
	address := models.Address{
		Id:           UpdateAddressRequest.Id,
		CustomerId:   customerId,
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
	Id int `validate:"min=0"`
}

type CreateAddressRequest struct {
	Lat          float64 `validate:"latitude,required"`
	Lon          float64 `validate:"longitude,required"`
	AddressLine1 string  `validate:"required,max=40"`
	AddressLine2 string  `validate:"max=40"`
	City         string  `validate:"required,max=40"`
	Country      string  `validate:"required,max=20"`
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

func CreateAddressRequestToAddress(createAddressRequest CreateAddressRequest, customerId int) models.Address {
	address := models.Address{
		CustomerId:   customerId,
		Lat:          createAddressRequest.Lat,
		Lon:          createAddressRequest.Lon,
		AddressLine1: createAddressRequest.AddressLine1,
		AddressLine2: createAddressRequest.AddressLine2,
		City:         createAddressRequest.City,
		Country:      createAddressRequest.Country,
	}

	return address
}

type GetAddressResponse struct {
	Id           int     `validate:"min=0"`
	Lat          float64 `validate:"latitude,required"`
	Lon          float64 `validate:"longitude,required"`
	AddressLine1 string  `validate:"required,max=40"`
	AddressLine2 string  `validate:"max=40"`
	City         string  `validate:"required,max=40"`
	Country      string  `validate:"required,max=20"`
}

func AddressToGetAddressResponse(address models.Address) GetAddressResponse {
	getAddressResponse := GetAddressResponse{
		Id:           address.Id,
		Lat:          address.Lat,
		Lon:          address.Lon,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		Country:      address.Country,
	}

	return getAddressResponse
}
