package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/VitoNaychev/food-app/auth"
	"github.com/VitoNaychev/food-app/customer-svc/models"
	"github.com/VitoNaychev/food-app/httperrors"
	"github.com/VitoNaychev/food-app/msgtypes"
	"github.com/VitoNaychev/food-app/validation"
	"github.com/golang-jwt/jwt/v5"
)

func (c *CustomerServer) AuthHandler(w http.ResponseWriter, r *http.Request) {
	authResponse := msgtypes.AuthResponse{Status: msgtypes.INVALID}

	if tokenHeader := r.Header.Get("Token"); tokenHeader == "" {
		handleAuthError(w, authResponse, msgtypes.MISSING_TOKEN)
		return
	}

	token, err := auth.VerifyJWT(r.Header["Token"][0], c.secretKey)
	if err != nil {
		handleAuthError(w, authResponse, msgtypes.INVALID)
		return
	}

	customerID, err := getCustomerIDFromToken(token)
	if err != nil {
		handleAuthError(w, authResponse, msgtypes.INVALID)
		return
	}

	_, err = c.store.GetCustomerByID(customerID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			handleAuthError(w, authResponse, msgtypes.NOT_FOUND)
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	authResponse.Status = msgtypes.OK
	authResponse.ID = customerID

	json.NewEncoder(w).Encode(authResponse)
}

func handleAuthError(w http.ResponseWriter, authResponse msgtypes.AuthResponse, status msgtypes.AuthStatus) {
	authResponse.Status = status
	json.NewEncoder(w).Encode(authResponse)
}

func getCustomerIDFromToken(token *jwt.Token) (int, error) {
	subject, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(subject)
}

func (c *CustomerServer) LoginHandler(w http.ResponseWriter, r *http.Request) {
	loginCustomerRequest, err := validation.ValidateBody[LoginCustomerRequest](r.Body)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	customer, err := c.store.GetCustomerByEmail(loginCustomerRequest.Email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			// wrap models.ErrNotFound in customer handlers error type?
			httperrors.WriteJSONError(w, http.StatusUnauthorized, ErrCustomerNotFound)
			return
		} else {
			httperrors.WriteJSONError(w, http.StatusInternalServerError, ErrDatabaseError)
			return
		}
	}

	if customer.Password != loginCustomerRequest.Password {
		httperrors.WriteJSONError(w, http.StatusUnauthorized, ErrInvalidCredentials)
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
		httperrors.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	customer = UpdateCustomerRequestToCustomer(updateCustomerRequest, id)

	err = c.store.UpdateCustomer(&customer)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusBadRequest, err)
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
		httperrors.WriteJSONError(w, http.StatusBadRequest, err)
		return
	}

	_, err = c.store.GetCustomerByEmail(createCustomerRequest.Email)
	if err == nil {
		httperrors.WriteJSONError(w, http.StatusBadRequest, ErrExistingCustomer)
		return
	} else if err != models.ErrNotFound {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, ErrDatabaseError)
		return
	}

	customer := CreateCustomerRequestToCustomer(createCustomerRequest)

	err = c.store.CreateCustomer(&customer)
	if err != nil {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, ErrDatabaseError)
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

func handleStoreError(w http.ResponseWriter, err error) {
	if errors.Is(err, models.ErrNotFound) {
		// wrap models.ErrNotFound in customer handlers error type?
		httperrors.WriteJSONError(w, http.StatusNotFound, ErrCustomerNotFound)
		return
	} else {
		httperrors.WriteJSONError(w, http.StatusInternalServerError, ErrDatabaseError)
		return
	}
}
