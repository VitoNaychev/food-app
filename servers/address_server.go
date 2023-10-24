package servers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type CustomerAddressStore interface {
	GetAddressByCustomerId(customerId int) (*GetAddressResponse, error)
}

type CustomerAddressServer struct {
	addressStore CustomerAddressStore
	secretKey    []byte
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
		authenticationMiddleware(c.GetAddressHandler, c.secretKey)(w, r)
	}
}

func (c *CustomerAddressServer) GetAddressHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.Header["Subject"][0])

	getAddressResponse, _ := c.addressStore.GetAddressByCustomerId(id)
	json.NewEncoder(w).Encode(getAddressResponse)
}
