package handlers

import "github.com/VitoNaychev/food-app/customer-svc/models"

type JWTResponse struct {
	Token string
}

type GetCustomerResponse struct {
	Id          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	Email       string `json:"email"`
}

func CustomerToGetCustomerResponse(customer models.Customer) GetCustomerResponse {
	getCustomerResponse := GetCustomerResponse{
		Id:          customer.Id,
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		PhoneNumber: customer.PhoneNumber,
		Email:       customer.Email,
	}

	return getCustomerResponse
}

type CreateCustomerRequest struct {
	FirstName   string `validate:"required,max=20"      json:"first_name"`
	LastName    string `validate:"required,max=20"      json:"last_name"`
	PhoneNumber string `validate:"required,phonenumber" json:"phone_number"`
	Email       string `validate:"required,email"       json:"email"`
	Password    string `validate:"required,max=72"      json:"password"`
}

func CustomerToCreateCustomerRequest(customer models.Customer) CreateCustomerRequest {
	createCustomerRequest := CreateCustomerRequest{
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Password:    customer.Password,
	}

	return createCustomerRequest
}

func CreateCustomerRequestToCustomer(createCustomerRequest CreateCustomerRequest) models.Customer {
	customer := models.Customer{
		Id:          0,
		FirstName:   createCustomerRequest.FirstName,
		LastName:    createCustomerRequest.LastName,
		PhoneNumber: createCustomerRequest.PhoneNumber,
		Email:       createCustomerRequest.Email,
		Password:    createCustomerRequest.Password,
	}
	return customer
}

type CreateCustomerResponse struct {
	JWT      JWTResponse
	Customer CustomerResponse
}

type CustomerResponse struct {
	Id          int    `validate:"min=1"                json:"id"`
	FirstName   string `validate:"required,max=20"      json:"first_name"`
	LastName    string `validate:"required,max=20"      json:"last_name"`
	PhoneNumber string `validate:"required,phonenumber" json:"phone_number"`
	Email       string `validate:"required,email"       json:"email"`
}

func CustomerToCustomerResponse(customer models.Customer) CustomerResponse {
	createCustomerResponse := CustomerResponse{
		Id:          customer.Id,
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		PhoneNumber: customer.PhoneNumber,
		Email:       customer.Email,
	}

	return createCustomerResponse
}

type LoginCustomerRequest struct {
	Email    string `validate:"required,email"       json:"email"`
	Password string `validate:"required,max=72"      json:"password"`
}

func CustomerToLoginCustomerRequest(customer models.Customer) LoginCustomerRequest {
	loginCustomerRequest := LoginCustomerRequest{
		Email:    customer.Email,
		Password: customer.Password,
	}

	return loginCustomerRequest
}

type UpdateCustomerRequest struct {
	FirstName   string `validate:"required,max=20"      json:"first_name"`
	LastName    string `validate:"required,max=20"      json:"last_name"`
	PhoneNumber string `validate:"required,phonenumber" json:"phone_number"`
	Email       string `validate:"required,email"       json:"email"`
	Password    string `validate:"required,max=72"      json:"password"`
}

func CustomerToUpdateCustomerRequest(customer models.Customer) UpdateCustomerRequest {
	updateCustomerRequest := UpdateCustomerRequest{
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Password:    customer.Password,
	}

	return updateCustomerRequest
}

func UpdateCustomerRequestToCustomer(updateCustomerRequest UpdateCustomerRequest, id int) models.Customer {
	customer := models.Customer{
		Id:          id,
		FirstName:   updateCustomerRequest.FirstName,
		LastName:    updateCustomerRequest.LastName,
		Email:       updateCustomerRequest.Email,
		PhoneNumber: updateCustomerRequest.PhoneNumber,
		Password:    updateCustomerRequest.Password,
	}

	return customer
}
