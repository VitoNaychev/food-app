package customer

import "github.com/VitoNaychev/bt-customer-svc/models/customer_store"

type GetCustomerResponse struct {
	Id          int
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
}

func CustomerToGetCustomerResponse(customer customer_store.Customer) GetCustomerResponse {
	getCustomerResponse := GetCustomerResponse{
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

func CustomerToCreateCustomerRequest(customer customer_store.Customer) CreateCustomerRequest {
	createCustomerRequest := CreateCustomerRequest{
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Password:    customer.Password,
	}

	return createCustomerRequest
}

func CreateCustomerRequestToCustomer(createCustomerRequest CreateCustomerRequest) customer_store.Customer {
	customer := customer_store.Customer{
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
	Id          int
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
}

func CustomerToCreateCustomerResponse(customer customer_store.Customer) CreateCustomerResponse {
	createCustomerResponse := CreateCustomerResponse{
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

func CustomerToLoginCustomerRequest(customer customer_store.Customer) LoginCustomerRequest {
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

func CustomerToUpdateCustomerRequest(customer customer_store.Customer) UpdateCustomerRequest {
	updateCustomerRequest := UpdateCustomerRequest{
		FirstName:   customer.FirstName,
		LastName:    customer.LastName,
		Email:       customer.Email,
		PhoneNumber: customer.PhoneNumber,
		Password:    customer.Password,
	}

	return updateCustomerRequest
}

func UpdateCustomerRequestToCustomer(updateCustomerRequest UpdateCustomerRequest, id int) customer_store.Customer {
	customer := customer_store.Customer{
		Id:          id,
		FirstName:   updateCustomerRequest.FirstName,
		LastName:    updateCustomerRequest.LastName,
		Email:       updateCustomerRequest.Email,
		PhoneNumber: updateCustomerRequest.PhoneNumber,
		Password:    updateCustomerRequest.Password,
	}

	return customer
}
