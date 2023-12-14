package handlers

import "github.com/VitoNaychev/food-app/courier-svc/models"

type LoginCourierRequest struct {
	Email    string `validate:"required,email,max=60"       json:"email"`
	Password string `valdiate:"required,max=72"             json:"password"`
}

type UpdateCourierRequest struct {
	FirstName   string `validate:"required,max=40"             json:"first_name"`
	LastName    string `validate:"required,max=40"             json:"last_name"`
	PhoneNumber string `validate:"required,phonenumber,max=20" json:"phone_number"`
	Email       string `validate:"required,email,max=60"       json:"email"`
	Password    string `valdiate:"required,max=72"             json:"password"`
	IBAN        string `validate:"required"                    json:"iban"`
}

func UpdateCourierRequestToCourier(request UpdateCourierRequest, id int) models.Courier {
	courier := models.Courier{
		ID:          id,
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
		Password:    request.Password,
		IBAN:        request.IBAN,
	}

	return courier
}

func CourierToUpdateCourierRequest(courier models.Courier) UpdateCourierRequest {
	updateCourierRequest := UpdateCourierRequest{
		FirstName:   courier.FirstName,
		LastName:    courier.LastName,
		PhoneNumber: courier.PhoneNumber,
		Email:       courier.Email,
		Password:    courier.Password,
		IBAN:        courier.IBAN,
	}

	return updateCourierRequest
}

type JWTResponse struct {
	Token string `validate:"required" json:"token"`
}

type CreateCourierResponse struct {
	JWT     JWTResponse     `validate:"required" json:"jwt"`
	Courier CourierResponse `validate:"required" json:"courier"`
}

type CourierResponse struct {
	ID          int    `validate:"min=1"                       json:"id"`
	FirstName   string `validate:"required,max=40"             json:"first_name"`
	LastName    string `validate:"required,max=40"             json:"last_name"`
	PhoneNumber string `validate:"required,phonenumber,max=20" json:"phone_number"`
	Email       string `validate:"required,email,max=60"       json:"email"`
	IBAN        string `validate:"required"                    json:"iban"`
}

func CourierToCourierResponse(courier models.Courier) CourierResponse {
	courierResponse := CourierResponse{
		ID:          courier.ID,
		FirstName:   courier.FirstName,
		LastName:    courier.LastName,
		PhoneNumber: courier.PhoneNumber,
		Email:       courier.Email,
		IBAN:        courier.IBAN,
	}

	return courierResponse
}

type CreateCourierRequest struct {
	FirstName   string `validate:"required,max=40"             json:"first_name"`
	LastName    string `validate:"required,max=40"             json:"last_name"`
	PhoneNumber string `validate:"required,phonenumber,max=20" json:"phone_number"`
	Email       string `validate:"required,email,max=60"       json:"email"`
	Password    string `valdiate:"required,max=72"             json:"password"`
	IBAN        string `validate:"required"                    json:"iban"`
}

func CourierToCreateCourierRequest(courier models.Courier) CreateCourierRequest {
	createCourierRequest := CreateCourierRequest{
		FirstName:   courier.FirstName,
		LastName:    courier.LastName,
		PhoneNumber: courier.PhoneNumber,
		Email:       courier.Email,
		Password:    courier.Password,
		IBAN:        courier.IBAN,
	}

	return createCourierRequest
}

func CreateCourierRequestToCourier(request CreateCourierRequest) models.Courier {
	courier := models.Courier{
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		PhoneNumber: request.PhoneNumber,
		Email:       request.Email,
		Password:    request.Password,
		IBAN:        request.IBAN,
	}

	return courier
}
