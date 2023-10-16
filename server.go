package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

type CustomerStore interface {
	GetCustomerInfo(id int) GetCustomerResponse
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
	if r.Header["Token"] != nil {
		if token, err := verifyJWT(r.Header["Token"][0]); err == nil {
			id := getIDFromToken(token)
			json.NewEncoder(w).Encode(c.store.GetCustomerInfo(id))
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}

	}
}

func verifyJWT(tokenString string) (*jwt.Token, error) {
	secretKey := []byte("mySecretKey")

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Invalid JWT token: %v", tokenString)
	}

	if token.Valid {
		return token, nil
	} else {
		return nil, fmt.Errorf("Expired JWT token")
	}
}

func getIDFromToken(token *jwt.Token) int {
	subject, _ := token.Claims.GetSubject()

	id, _ := strconv.Atoi(subject)

	return id
}
