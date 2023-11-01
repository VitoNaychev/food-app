package customer_store

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/VitoNaychev/bt-customer-svc/handlers"
	"github.com/VitoNaychev/bt-customer-svc/handlers/customer"
	"github.com/VitoNaychev/bt-customer-svc/models/customer_store"
	"github.com/VitoNaychev/bt-customer-svc/testdata"
	"github.com/VitoNaychev/bt-customer-svc/testutil"
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
		postgres.WithInitScripts(filepath.Join("testdata", "init-db.sql")),
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

func TestCreateCustomerAndRetrievThem(t *testing.T) {
	connStr := SetupDatabaseContainer(t)

	store, err := customer_store.NewPostgresCustomerStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	server := customer.NewCustomerServer(testEnv.SecretKey, testEnv.ExpiresAt, &store)

	request := customer.NewCreateCustomerRequest(testdata.PeterCustomer)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Header()["Token"] == nil {
		t.Fatalf("server didn't return JWT")
	}

	request = customer.NewGetCustomerRequest(response.Header()["Token"][0])
	response = httptest.NewRecorder()

	server.ServeHTTP(response, request)

	// want := customerToGetCustomerResponse(testdata.PeterCustomer)

	var got customer.GetCustomerResponse
	json.NewDecoder(response.Body).Decode(&got)

	testutil.AssertStatus(t, response.Code, http.StatusOK)
	// assertGetCustomerResponse(t, got, want)
}
