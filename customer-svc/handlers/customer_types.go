package handlers

import "github.com/VitoNaychev/bt-customer-svc/models"

type AuthStatus int

const (
	MISSING_TOKEN AuthStatus = iota
	INVALID
	NOT_FOUND
	OK
)

type AuthResponse struct {
	Status AuthStatus
	ID     int
}

type JWTResponse struct {
	Token string
}

type GetCustomerResponse struct {
	Id          int
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
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
	FirstName   string `validate:"required,max=20"`
	LastName    string `validate:"required,max=20"`
	PhoneNumber string `validate:"required,phonenumber"`
	Email       string `validate:"required,email"`
	Password    string `validate:"required,max=72"`
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
	Id          int
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
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
	Email    string `validate:"required,email"`
	Password string `validate:"required,max=72"`
}

func CustomerToLoginCustomerRequest(customer models.Customer) LoginCustomerRequest {
	loginCustomerRequest := LoginCustomerRequest{
		Email:    customer.Email,
		Password: customer.Password,
	}

	return loginCustomerRequest
}

type UpdateCustomerRequest struct {
	FirstName   string `validate:"required,max=20"`
	LastName    string `validate:"required,max=20"`
	PhoneNumber string `validate:"phonenumber,required"`
	Email       string `validate:"required,email"`
	Password    string `validate:"required,max=72"`
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
