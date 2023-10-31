package bt_customer_svc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupDatabaseContainer(t testing.TB) string {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("testdata", "init-db.sql")),
		postgres.WithDatabase(testEnv.dbname),
		postgres.WithUsername(testEnv.dbuser),
		postgres.WithPassword(testEnv.dbpass),
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

	store, err := NewPostgresCustomerStore(context.Background(), connStr)
	if err != nil {
		t.Fatal(err)
	}

	server := NewCustomerServer(testEnv.secretKey, testEnv.expiresAt, &store)

	request := newCreateCustomerRequest(peterCustomer)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	if response.Header()["Token"] == nil {
		t.Fatalf("server didn't return JWT")
	}

	request = newGetCustomerRequest(response.Header()["Token"][0])
	response = httptest.NewRecorder()

	server.ServeHTTP(response, request)

	want := customerToGetCustomerResponse(peterCustomer)

	var got GetCustomerResponse
	json.NewDecoder(response.Body).Decode(&got)

	assertStatus(t, response.Code, http.StatusOK)
	assertGetCustomerResponse(t, got, want)
}
