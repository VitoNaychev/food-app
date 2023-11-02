package customer

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/handlers/auth"
	"github.com/VitoNaychev/bt-customer-svc/handlers/validation"
)

func (c *CustomerServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginCustomerRequest LoginCustomerRequest
	err := validation.ValidateBody(r.Body, &loginCustomerRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: err.Error()})
		return
	}

	customer, err := c.store.GetCustomerByEmail(loginCustomerRequest.Email)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingCustomer.Error()})
		return
	}

	if customer.Password != loginCustomerRequest.Password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrInvalidCredentials.Error()})
		return
	}

	loginJWT, _ := auth.GenerateJWT(c.secretKey, c.expiresAt, customer.Id)

	w.WriteHeader(http.StatusAccepted)
	w.Header().Add("Token", loginJWT)
}

func (c *CustomerServer) updateCustomer(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Header["Subject"][0])
	customer, _ := c.store.GetCustomerByID(id)

	var updateCustomerRequest UpdateCustomerRequest
	err := validation.ValidateBody(r.Body, &updateCustomerRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: err.Error()})
		return
	}

	customer = UpdateCustomerRequestToCustomer(updateCustomerRequest, id)

	err = c.store.UpdateCustomer(&customer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrDatabaseError.Error()})
	}

	json.NewEncoder(w).Encode(CustomerToCustomerResponse(customer))
}

func (c *CustomerServer) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Header["Subject"][0])
	err := c.store.DeleteCustomer(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingCustomer.Error()})
	}
}

func (c *CustomerServer) createCustomer(w http.ResponseWriter, r *http.Request) {
	var createCustomerRequest CreateCustomerRequest
	err := validation.ValidateBody(r.Body, &createCustomerRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: err.Error()})
		return
	}

	_, err = c.store.GetCustomerByEmail(createCustomerRequest.Email)
	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrExistingUser.Error()})
		return
	}

	customer := CreateCustomerRequestToCustomer(createCustomerRequest)

	err = c.store.CreateCustomer(&customer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrDatabaseError.Error()})
	}

	customerJWT, _ := auth.GenerateJWT(c.secretKey, c.expiresAt, customer.Id)

	w.WriteHeader(http.StatusAccepted)
	w.Header().Add("Token", customerJWT)
	json.NewEncoder(w).Encode(CustomerToCustomerResponse(customer))
}

func (c *CustomerServer) getCustomer(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Header["Subject"][0])
	customer, err := c.store.GetCustomerByID(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: handlers.ErrMissingCustomer.Error()})
		return
	}

	getCustomerResponse := CustomerToGetCustomerResponse(customer)
	json.NewEncoder(w).Encode(getCustomerResponse)

}
