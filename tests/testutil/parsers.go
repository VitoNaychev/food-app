package testutil

import (
	"encoding/json"
	"io"
	"testing"

	"github.com/VitoNaychev/bt-customer-svc/handlers/address"
	"github.com/VitoNaychev/bt-customer-svc/handlers/customer"
	"github.com/VitoNaychev/bt-customer-svc/models"
)

func ParseCustomerResponse(t testing.TB, r io.Reader) (customerResponse customer.CustomerResponse) {
	t.Helper()

	json.NewDecoder(r).Decode(&customerResponse)
	return
}

func ParseAddressResponse(t testing.TB, r io.Reader) (addressResponse models.Address) {
	t.Helper()

	json.NewDecoder(r).Decode(&addressResponse)
	return
}

func ParseGetAddressResponse(t testing.TB, r io.Reader) (getAddressResponse []address.GetAddressResponse) {
	t.Helper()

	json.NewDecoder(r).Decode(&getAddressResponse)
	return
}
