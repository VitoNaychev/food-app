package stubs

import (
	"reflect"
	"testing"

	"github.com/VitoNaychev/food-app/customer-svc/models"
)

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

func AssertCreatedCustomer(t testing.TB, store *StubCustomerStore, customer models.Customer) {
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
