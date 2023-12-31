package pgconfig

import (
	"fmt"
	"strings"

	"github.com/VitoNaychev/food-app/appenv"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Options  []string
}

func (c *Config) GetConnectionString() string {
	optionsStr := strings.Join(c.Options, "&")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		c.User, c.Password, c.Host, c.Port, c.Database, optionsStr)
}

func GetConfigFromEnv(env appenv.Enviornment) Config {
	config := Config{
		Host:     env.Dbhost,
		Port:     env.Dbport,
		User:     env.Dbuser,
		Password: env.Dbpass,
		Database: env.Dbname,
	}

	return config
}
