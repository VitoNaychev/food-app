package integrationtest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/handlers/customer"
	"github.com/VitoNaychev/bt-customer-svc/models"
	"github.com/VitoNaychev/bt-customer-svc/tests/testdata"
	"github.com/VitoNaychev/bt-customer-svc/tests/testutil"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testEnv handlers.TestEnv

func TestMain(m *testing.M) {
	testEnv = handlers.LoadTestEnviornment()

	code := m.Run()
	os.Exit(code)
}

func SetupDatabaseContainer(t testing.TB) string {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("..", "testdata", "init-db.sql")),
		postgres.WithDatabase(testEnv.Dbname),
		postgres.WithUsername(testEnv.Dbuser),
		postgres.WithPassword(testEnv.Dbpass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second),
		))

	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("couldn't get connection string: %v", err)
	}

	return connStr
}

func TestCustomerServerOperations(t *testing.T) {
	connStr := SetupDatabaseContainer(t)

	store, err := models.NewPgCustomerStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	server := customer.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, &store)

	peterJWT := createNewCustomer(server, testdata.PeterCustomer)
	aliceJWT := createNewCustomer(server, testdata.AliceCustomer)

	t.Run("retrieving customer", func(t *testing.T) {
		request := customer.NewGetCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := customer.CustomerToCustomerResponse(testdata.PeterCustomer)
		var got customer.CustomerResponse
		json.NewDecoder(response.Body).Decode(&got)
		AssertCustomerResponse(t, got, want)
	})

	t.Run("updating customer", func(t *testing.T) {
		updateCustomer := testdata.AliceCustomer
		updateCustomer.LastName = "Roper"

		request := customer.NewUpdateCustomerRequest(updateCustomer, aliceJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		want := customer.CustomerToCustomerResponse(updateCustomer)
		var got customer.CustomerResponse
		json.NewDecoder(response.Body).Decode(&got)
		AssertCustomerResponse(t, got, want)
	})

	t.Run("deleting customer", func(t *testing.T) {
		request := customer.NewDeleteCustomerRequest(peterJWT)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusOK)

		request = customer.NewGetCustomerRequest(peterJWT)
		response = httptest.NewRecorder()

		server.ServeHTTP(response, request)

		testutil.AssertStatus(t, response.Code, http.StatusNotFound)
	})
}

func createNewCustomer(server http.Handler, c models.Customer) string {
	request := customer.NewCreateCustomerRequest(c)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	return response.Header()["Token"][0]
}

func AssertCustomerResponse(t testing.TB, got, want customer.CustomerResponse) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got response %v want %v", got, want)
	}
}
