package tests

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type TestEnv struct {
	SecretKey []byte
	ExpiresAt time.Duration

	Dbuser string
	Dbpass string
	Dbname string
}

func LoadTestEnviornment() TestEnv {
	godotenv.Load("../test.env")

	testEnv := TestEnv{}
	testEnv.SecretKey = []byte(os.Getenv("SECRET"))
	testEnv.ExpiresAt = time.Second

	testEnv.Dbuser = os.Getenv("DBUSER")
	testEnv.Dbpass = os.Getenv("DBPASS")
	testEnv.Dbname = os.Getenv("DBNAME")

	return testEnv
}
