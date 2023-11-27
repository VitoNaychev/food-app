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

	env := Enviornment{}
	env.SecretKey = []byte(os.Getenv("SECRET"))
	env.ExpiresAt = time.Second

	env.Dbuser = os.Getenv("DBUSER")
	env.Dbpass = os.Getenv("DBPASS")
	env.Dbname = os.Getenv("DBNAME")

	return env
}
