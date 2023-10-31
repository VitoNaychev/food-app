package bt_customer_svc

import (
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

type TestEnv struct {
	secretKey []byte
	expiresAt time.Duration

	dbuser string
	dbpass string
	dbname string
}

var testEnv TestEnv

func TestMain(m *testing.M) {
	godotenv.Load("test.env")

	testEnv.secretKey = []byte(os.Getenv("SECRET"))
	testEnv.expiresAt = time.Second

	testEnv.dbuser = os.Getenv("DBUSER")
	testEnv.dbpass = os.Getenv("DBPASS")
	testEnv.dbname = os.Getenv("DBNAME")

	code := m.Run()

	os.Exit(code)
}
