package integration

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/VitoNaychev/food-app/loadenv"
	"github.com/VitoNaychev/food-app/testutil"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func NewDatabaseContainer(t testing.TB) string {
	keys := []string{"DBUSER", "DBPASS", "DBNAME"}

	env, err := loadenv.LoadEnviornment("../config/test.env", keys)
	if err != nil {
		testutil.HandleLoadEnviornmentError(err)
		t.Fatal()
	}

	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithInitScripts(filepath.Join("..", "sql-scripts", "init.sql")),
		postgres.WithDatabase(env.Dbname),
		postgres.WithUsername(env.Dbuser),
		postgres.WithPassword(env.Dbpass),
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
