package customer

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/handlers/auth"
	"github.com/VitoNaychev/bt-customer-svc/models"
	"github.com/VitoNaychev/validation"
)

func (c *CustomerServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginCustomerRequest, err := validation.ValidateBody[LoginCustomerRequest](r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	customer, err := c.store.GetCustomerByEmail(loginCustomerRequest.Email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			// wrap models.ErrNotFound in customer handlers error type?
			writeJSONError(w, http.StatusUnauthorized, handlers.ErrCustomerNotFound)
			return
		} else {
			writeJSONError(w, http.StatusInternalServerError, handlers.ErrDatabaseError)
			return
		}
	}

	if customer.Password != loginCustomerRequest.Password {
		writeJSONError(w, http.StatusUnauthorized, handlers.ErrInvalidCredentials)
		return
	}

	loginJWT, _ := auth.GenerateJWT(c.secretKey, c.expiresAt, customer.Id)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(JWTResponse{Token: loginJWT})
}

func (c *CustomerServer) updateCustomer(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Header["Subject"][0])
	customer, err := c.store.GetCustomerByID(id)
	if err != nil {
		handleStoreError(w, err)
	}

	updateCustomerRequest, err := validation.ValidateBody[UpdateCustomerRequest](r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	customer = UpdateCustomerRequestToCustomer(updateCustomerRequest, id)

	err = c.store.UpdateCustomer(&customer)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err)
	}

	json.NewEncoder(w).Encode(CustomerToCustomerResponse(customer))
}

func (c *CustomerServer) deleteCustomer(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Header["Subject"][0])
	err := c.store.DeleteCustomer(id)

	if err != nil {
		handleStoreError(w, err)
	}
}

func (c *CustomerServer) createCustomer(w http.ResponseWriter, r *http.Request) {
	createCustomerRequest, err := validation.ValidateBody[CreateCustomerRequest](r.Body)
	if err != nil {
		// Add error wrapping for validation errors
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	_, err = c.store.GetCustomerByEmail(createCustomerRequest.Email)
	if err == nil {
		writeJSONError(w, http.StatusBadRequest, handlers.ErrExistingCustomer)
		return
	} else if err != models.ErrNotFound {
		writeJSONError(w, http.StatusInternalServerError, handlers.ErrDatabaseError)
		return
	}

	customer := CreateCustomerRequestToCustomer(createCustomerRequest)

	err = c.store.CreateCustomer(&customer)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, handlers.ErrDatabaseError)
		return
	}

	customerJWT, _ := auth.GenerateJWT(c.secretKey, c.expiresAt, customer.Id)

	w.WriteHeader(http.StatusAccepted)

	createCustomerResponse := CreateCustomerResponse{
		JWT:      JWTResponse{Token: customerJWT},
		Customer: CustomerToCustomerResponse(customer),
	}
	json.NewEncoder(w).Encode(createCustomerResponse)
}

func (c *CustomerServer) getCustomer(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Header["Subject"][0])
	customer, err := c.store.GetCustomerByID(id)

	if err != nil {
		handleStoreError(w, err)
		return
	}

	getCustomerResponse := CustomerToGetCustomerResponse(customer)
	json.NewEncoder(w).Encode(getCustomerResponse)
}

func writeJSONError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(handlers.ErrorResponse{Message: err.Error()})
}

func handleStoreError(w http.ResponseWriter, err error) {
	if errors.Is(err, models.ErrNotFound) {
		// wrap models.ErrNotFound in customer handlers error type?
		writeJSONError(w, http.StatusNotFound, handlers.ErrCustomerNotFound)
		return
	} else {
		writeJSONError(w, http.StatusInternalServerError, handlers.ErrDatabaseError)
		return
	}
}
