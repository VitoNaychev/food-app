package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type CustomerStore interface {
	GetCustomerInfo(id int) (*GetCustomerResponse, error)
}

type CustomerServer struct {
	store CustomerStore
}

type GetCustomerResponse struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
}

type CreateCustomerRequest struct {
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Password    string
}

type ErrorResponse struct {
	Message string
}

func (c *CustomerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		storeCustomer(w, r)
	case http.MethodGet:
		c.getCustomer(w, r)
	}
}

func storeCustomer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func (c *CustomerServer) getCustomer(w http.ResponseWriter, r *http.Request) {
	if token, err := verifyJWT(r.Header); err == nil {
		id := getIDFromToken(token)
		customerResponse, err := c.store.GetCustomerInfo(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
		} else {
			json.NewEncoder(w).Encode(customerResponse)
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
	}
}

func verifyJWT(header http.Header) (*jwt.Token, error) {
	secretKey := []byte("mySecretKey")

	if header["Token"] == nil {
		return nil, errors.New("token is missing")
	}

	token, err := jwt.Parse(header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {
		return nil, err
	}

	return token, nil
}

func getIDFromToken(token *jwt.Token) int {
	subject, _ := token.Claims.GetSubject()

	id, _ := strconv.Atoi(subject)

	return id
}
