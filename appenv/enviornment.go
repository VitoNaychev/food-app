package appenv

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var ErrInvalidDuration = errors.New("duration string is invalid")
var ErrUnsupportedVariable = errors.New("trying to load an unsupported enviornment variable")

type Enviornment struct {
	SecretKey []byte
	ExpiresAt time.Duration

	Dbuser string
	Dbpass string
	Dbname string

	KafkaBrokers []string
}

func LoadEnviornment(file string, keys []string) (Enviornment, error) {
	godotenv.Load(file)

	env := Enviornment{}

	var err error
	for _, key := range keys {
		switch key {
		case "SECRET":
			var strSecretKey string
			strSecretKey, err = getRequiredEnv("SECRET")

			env.SecretKey = []byte(strSecretKey)
		case "EXPIRES_AT":
			var expiresAtStr string
			expiresAtStr, err = getRequiredEnv("EXPIRES_AT")

			if err == nil {
				env.ExpiresAt, err = parseDuration(expiresAtStr)
			}
		case "DBUSER":
			env.Dbuser, err = getRequiredEnv("DBUSER")
		case "DBPASS":
			env.Dbpass, err = getRequiredEnv("DBPASS")
		case "DBNAME":
			env.Dbname, err = getRequiredEnv("DBNAME")
		default:
			return Enviornment{}, ErrUnsupportedVariable
		}

		if err != nil {
			return Enviornment{}, err
		}
	}

	return env, nil
}

func getRequiredEnv(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("enviornment variable %v is not set", name)
	}
	return value, nil
}

func parseDuration(durationStr string) (time.Duration, error) {
	parsedTime, err := time.Parse("15:04:05", durationStr)
	if err != nil {
		return 0, ErrInvalidDuration
	}

	duration := parsedTime.Sub(time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC))

	return duration, nil
}
