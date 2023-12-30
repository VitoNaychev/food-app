package pgconfig

import (
	"fmt"

	"github.com/VitoNaychev/food-app/appenv"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

func (c *Config) GetConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		c.User, c.Password, c.Host, c.Port, c.Database)
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
