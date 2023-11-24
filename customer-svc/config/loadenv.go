package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Enviornment struct {
	SecretKey []byte
	ExpiresAt time.Duration

	Dbuser string
	Dbpass string
	Dbname string
}

func LoadEnviornment(file string) Enviornment {
	godotenv.Load(file)

	testEnv := Enviornment{}
	testEnv.SecretKey = []byte(os.Getenv("SECRET"))
	testEnv.ExpiresAt = time.Second

	testEnv.Dbuser = os.Getenv("DBUSER")
	testEnv.Dbpass = os.Getenv("DBPASS")
	testEnv.Dbname = os.Getenv("DBNAME")

	return testEnv
}
