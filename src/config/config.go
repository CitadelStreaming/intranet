package config

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	ENV_SECRET_LOCATION             = "SECRET_LOCATION"
	ENV_SECRET_LOCATION_ENVIRONMENT = "environment"
	ENV_SECRET_LOCATION_SECRETS     = "secrets"

	ENV_DATABASE_HOST = "DB_HOST"
	ENV_DATABASE_USER = "DB_USER"
	ENV_DATABASE_PASS = "DB_PASS"
	ENV_DATABASE_PORT = "DB_PORT"
	ENV_DATABASE_NAME = "DB_NAME"

	ENV_MIGRATIONS_PATH = "MIGRATIONS"
)

type Config struct {
	DbHost string
	DbUser string
	DbPass string
	DbPort uint16
	DbName string

	MigrationsPath string
}

/*
Load a configuration from environment variables, providing default values
*/
func LoadConfig() Config {
	return Config{
		DbHost: getEnvStringWithDefault(ENV_DATABASE_HOST, "localhost"),
		DbUser: getEnvStringWithDefault(ENV_DATABASE_USER, ""),
		DbPass: getEnvStringWithDefault(ENV_DATABASE_PASS, ""),
		DbPort: getEnvUint16WithDefault(ENV_DATABASE_PORT, 3306),
		DbName: getEnvStringWithDefault(ENV_DATABASE_NAME, ""),

		MigrationsPath: getEnvStringWithDefault(ENV_MIGRATIONS_PATH, "/var/migrations"),
	}
}

func getEnvStringWithDefault(key string, defaultValue string) string {
	if val, found := os.LookupEnv(key); found {
		return val
	}
	return defaultValue
}

func getEnvUint16WithDefault(key string, defaultValue uint16) uint16 {
	if val, found := os.LookupEnv(key); found {
		ret, err := strconv.ParseUint(val, 10, 16)
		if err == nil {
			return uint16(ret)
		} else {
			logrus.Warn("Unable to parse uint16 value, returning default.", err.Error())
		}
	}
	return defaultValue
}
