package address

import "github.com/VitoNaychev/bt-customer-svc/models/address_store"

type UpdateAddressRequest struct {
	Id           int     `validate:"min=0"`
	Lat          float64 `validate:"latitude,required"`
	Lon          float64 `validate:"longitude,required"`
	AddressLine1 string  `validate:"required,max=40"`
	AddressLine2 string  `validate:"max=40"`
	City         string  `validate:"required,max=40"`
	Country      string  `validate:"required,max=20"`
}

func AddressToUpdateAddressRequest(address address_store.Address) UpdateAddressRequest {
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

func UpdateAddressRequestToAddress(UpdateAddressRequest UpdateAddressRequest, customerId int) address_store.Address {
	address := address_store.Address{
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

type AddAddressRequest struct {
	Lat          float64 `validate:"latitude,required"`
	Lon          float64 `validate:"longitude,required"`
	AddressLine1 string  `validate:"required,max=40"`
	AddressLine2 string  `validate:"max=40"`
	City         string  `validate:"required,max=40"`
	Country      string  `validate:"required,max=20"`
}

func AddressToAddAddressRequest(address address_store.Address) AddAddressRequest {
	addAddressRequest := AddAddressRequest{
		Lat:          address.Lat,
		Lon:          address.Lon,
		AddressLine1: address.AddressLine1,
		AddressLine2: address.AddressLine2,
		City:         address.City,
		Country:      address.Country,
	}

	return addAddressRequest
}

func AddAddressRequestToAddress(addAddressRequest AddAddressRequest, customerId int) address_store.Address {
	address := address_store.Address{
		CustomerId:   customerId,
		Lat:          addAddressRequest.Lat,
		Lon:          addAddressRequest.Lon,
		AddressLine1: addAddressRequest.AddressLine1,
		AddressLine2: addAddressRequest.AddressLine2,
		City:         addAddressRequest.City,
		Country:      addAddressRequest.Country,
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

func AddressToGetAddressResponse(address address_store.Address) GetAddressResponse {
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
