package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

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

func CustomerServer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		storeCustomer(w, r)
	case http.MethodGet:
		getCustomer(w, r)
	}
}

func storeCustomer(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func getCustomer(w http.ResponseWriter, r *http.Request) {
	secretKey := []byte("mySecretKey")

	peterResponse := GetCustomerResponse{
		FirstName:   "Peter",
		LastName:    "Smith",
		PhoneNumber: "+359 88 576 5981",
		Email:       "petesmith@gmail.com",
	}

	aliceResponse := GetCustomerResponse{
		FirstName:   "Alice",
		LastName:    "Johnson",
		PhoneNumber: "+359 88 444 2222",
		Email:       "alicejohn@gmail.com",
	}

	var response GetCustomerResponse

	if r.Header["Token"] != nil {
		token, _ := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return "", fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return secretKey, nil
		})

		subject, _ := token.Claims.GetSubject()

		switch subject {
		case "0":
			response = peterResponse
		case "1":
			response = aliceResponse
		}
	}

	json.NewEncoder(w).Encode(response)
}
