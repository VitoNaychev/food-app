package testutil

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/VitoNaychev/food-app/customer-svc/handlers"
	"github.com/VitoNaychev/food-app/customer-svc/models"
)

func ParseCreateCustomerResponse(t testing.TB, r io.Reader) (createCustomerResponse handlers.CreateCustomerResponse) {
	t.Helper()

	json.NewDecoder(r).Decode(&createCustomerResponse)
	return
}

func ParseCustomerResponse(t testing.TB, r io.Reader) (customerResponse handlers.CustomerResponse) {
	t.Helper()

	json.NewDecoder(r).Decode(&customerResponse)
	return
}

func ParseAddressResponse(t testing.TB, r io.Reader) (addressResponse models.Address) {
	t.Helper()

	json.NewDecoder(r).Decode(&addressResponse)
	return
}

func ParseGetAddressResponse(t testing.TB, r io.Reader) (getAddressResponse []handlers.GetAddressResponse) {
	t.Helper()

	json.NewDecoder(r).Decode(&getAddressResponse)
	return
}
