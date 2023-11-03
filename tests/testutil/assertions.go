package testutil

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"testing"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/handlers/address"
	"github.com/VitoNaychev/bt-customer-svc/handlers/customer"
	"github.com/VitoNaychev/bt-customer-svc/models"
	"github.com/golang-jwt/jwt/v5"
)

func AssertAddresses(t testing.TB, got, want []address.GetAddressResponse) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func AssertCustomerResponse(t testing.TB, got, want customer.CustomerResponse) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got response %v want %v", got, want)
	}
}

func AssertErrorResponse(t testing.TB, body io.Reader, expetedError error) {
	t.Helper()

	var errorResponse handlers.ErrorResponse
	json.NewDecoder(body).Decode(&errorResponse)

	if errorResponse.Message != expetedError.Error() {
		t.Errorf("got error %q want %q", errorResponse.Message, expetedError.Error())
	}
}

func AssertStatus(t testing.TB, got, want int) {
	t.Helper()

	if got != want {
		t.Fatalf("got status %d, want %d", got, want)
	}
}

func AssertValidResponse(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("couldn't parse response body, got error %v", err)
	}
}

func AssertJWT(t testing.TB, header http.Header, secretKey []byte, wantId int) {
	t.Helper()

	if header["Token"] == nil {
		t.Fatalf("missing JWT in header")
	}

	token, err := jwt.Parse(header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))

	if err != nil {
		t.Fatalf("error verifying JWT: %v", err)
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		t.Fatalf("did not get subject in JWT, expected %v", wantId)
	}

	gotId, err := strconv.Atoi(subject)
	if err != nil {
		t.Fatalf("did not get customer ID for subject, got %v", subject)
	}

	if gotId != wantId {
		t.Errorf("got customer id %v want %v", subject, wantId)
	}
}

func AssertUpdatedCustomer(t testing.TB, store *StubCustomerStore, customer models.Customer) {
	t.Helper()

	if len(store.updateCalls) != 1 {
		t.Fatalf("got %d calls to UpdateCustomer expected %d", len(store.updateCalls), 1)
	}

	if !reflect.DeepEqual(store.updateCalls[0], customer) {
		t.Errorf("did not update correct customer got %v want %v", store.updateCalls[0], customer)
	}
}

func AssertDeletedCustomer(t testing.TB, store *StubCustomerStore, customer models.Customer) {
	t.Helper()

	if len(store.deleteCalls) != 1 {
		t.Fatalf("got %d calls to DeleteCustomer expected %d", len(store.deleteCalls), 1)
	}

	if store.deleteCalls[0] != customer.Id {
		t.Errorf("did not delete correct customer got %d want %d", store.deleteCalls[0], customer.Id)
	}
}

func AssertStoredCustomer(t testing.TB, store *StubCustomerStore, customer models.Customer) {
	t.Helper()

	if len(store.storeCalls) != 1 {
		t.Fatalf("got %d calls to StoreAddress expected %d", len(store.storeCalls), 1)
	}

	// Copy the address ID from the dummy data to the stored data. This is
	// done because the tests use stubs for store implementations and
	// dummy user data for test cases, so there is going to be a mismatch
	// of the assigned IDs between the two.
	store.storeCalls[0].Id = customer.Id
	if !reflect.DeepEqual(store.storeCalls[0], customer) {
		t.Errorf("did not store correct customer got %d want %d", store.storeCalls[0].Id, customer.Id)
	}
}

func AssertUpdatedAddress(t testing.TB, store *StubAddressStore, address models.Address) {
	t.Helper()

	if len(store.updateCalls) != 1 {
		t.Fatalf("got %d calls to UpdateAddress expected %d", len(store.updateCalls), 1)
	}

	if !reflect.DeepEqual(store.updateCalls[0], address) {
		t.Errorf("did not update correct address got %v want %v", store.updateCalls[0], address)
	}
}

func AssertDeletedAddress(t testing.TB, store *StubAddressStore, address models.Address) {
	t.Helper()

	if len(store.deleteCalls) != 1 {
		t.Fatalf("got %d calls to DeleteAddress expected %d", len(store.updateCalls), 1)
	}

	if store.deleteCalls[0] != address.Id {
		t.Errorf("did not delete correct address got %d want %d", store.deleteCalls[0], address.Id)
	}
}

func AssertStoredAddress(t testing.TB, store *StubAddressStore, address models.Address) {
	t.Helper()

	if len(store.storeCalls) != 1 {
		t.Fatalf("got %d calls to StoreAddress expected %d", len(store.storeCalls), 1)
	}

	// Copy the address ID from the dummy data to the stored data. This is
	// done because the tests use stubs for store implementations and
	// dummy user data for test cases, so there is going to be a mismatch
	// of the assigned IDs between the two.
	store.storeCalls[0].Id = address.Id
	if !reflect.DeepEqual(store.storeCalls[0], address) {
		t.Errorf("did not store correct customer got %d want %d", store.storeCalls[0].Id, address.Id)
	}
}
