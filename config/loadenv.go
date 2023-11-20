package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Enviornment struct {
	Dbuser string
	Dbpass string
	Dbname string
}

func LoadEnviornment(file string) Enviornment {
	godotenv.Load(file)

	env := Enviornment{
		Dbuser: os.Getenv("DBUSER"),
		Dbpass: os.Getenv("DBPASS"),
		Dbname: os.Getenv("DBNAME"),
	}

	return env
}
