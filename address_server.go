package bt_customer_svc

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type CustomerAddressStore interface {
	GetAddressesByCustomerId(customerId int) ([]GetAddressResponse, error)
}

type CustomerAddressServer struct {
	addressStore  CustomerAddressStore
	customerStore CustomerStore
	secretKey     []byte
}

type GetAddressResponse struct {
	Lat          float64
	Lon          float64
	AddressLine1 string
	AddressLine2 string
	City         string
	Country      string
}

func (c *CustomerAddressServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		AuthenticationMiddleware(c.GetAddressHandler, c.secretKey)(w, r)
	}
}

func (c *CustomerAddressServer) GetAddressHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Header["Subject"][0])

	_, err := c.customerStore.GetCustomerById(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Message: ErrMissingCustomer.Error()})
		return
	}

	getAddressResponse, _ := c.addressStore.GetAddressesByCustomerId(id)
	json.NewEncoder(w).Encode(getAddressResponse)
}
