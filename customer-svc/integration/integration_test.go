package integrationtest

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/VitoNaychev/bt-customer-svc/config"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testEnv config.Enviornment

func TestMain(m *testing.M) {
	testEnv = config.LoadEnviornment("../config/test.env")
	fmt.Println("+++ENV: ", testEnv)

	code := m.Run()
	os.Exit(code)
}

func SetupDatabaseContainer(t testing.TB) string {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("..", "sql-scripts", "init.sql")),
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
