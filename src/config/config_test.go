package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"citadel_intranet/src/config"
)

func TestLoadConfigEmpty(t *testing.T) {
	assert := assert.New(t)
	cfg := config.LoadConfig()

	assert.Equal("localhost", cfg.DbHost)
	assert.Equal("", cfg.DbUser)
	assert.Equal("", cfg.DbPass)
	assert.Equal(uint16(3306), cfg.DbPort)
	assert.Equal("", cfg.DbName)

	assert.Equal("localhost", cfg.ServerHost)
	assert.Equal(uint16(8080), cfg.ServerPort)
	assert.Equal("/var/www", cfg.ServerFilePath)

	assert.Equal("/var/migrations", cfg.MigrationsPath)
}

func TestLoadConfigSetValues(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(os.Setenv(config.ENV_DATABASE_HOST, "database.local"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_USER, "bobby"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_PASS, "tables"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_PORT, "23306"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_NAME, "db1"))

	assert.Nil(os.Setenv(config.ENV_SERVER_HOST, "webserver.local"))
	assert.Nil(os.Setenv(config.ENV_SERVER_PORT, "80"))
	assert.Nil(os.Setenv(config.ENV_SERVER_PATH, "/var/www/site1"))

	assert.Nil(os.Setenv(config.ENV_MIGRATIONS_PATH, "/opt/citadel/migrations"))

	cfg := config.LoadConfig()

	assert.Equal("database.local", cfg.DbHost)
	assert.Equal("bobby", cfg.DbUser)
	assert.Equal("tables", cfg.DbPass)
	assert.Equal(uint16(23306), cfg.DbPort)
	assert.Equal("db1", cfg.DbName)

	assert.Equal("webserver.local", cfg.ServerHost)
	assert.Equal(uint16(80), cfg.ServerPort)
	assert.Equal("/var/www/site1", cfg.ServerFilePath)

	assert.Equal("/opt/citadel/migrations", cfg.MigrationsPath)
}

func TestLoadConfigSetValuesInvalidPort(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(os.Setenv(config.ENV_DATABASE_HOST, "database.local"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_USER, "bobby"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_PASS, "tables"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_PORT, "NotANumber"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_NAME, "db1"))

	assert.Nil(os.Setenv(config.ENV_SERVER_HOST, "webserver.local"))
	assert.Nil(os.Setenv(config.ENV_SERVER_PORT, "waggles"))
	assert.Nil(os.Setenv(config.ENV_SERVER_PATH, "/var/www/site1"))

	assert.Nil(os.Setenv(config.ENV_MIGRATIONS_PATH, "/var/migrations"))

	cfg := config.LoadConfig()

	assert.Equal("database.local", cfg.DbHost)
	assert.Equal("bobby", cfg.DbUser)
	assert.Equal("tables", cfg.DbPass)
	assert.Equal(uint16(3306), cfg.DbPort)
	assert.Equal("db1", cfg.DbName)

	assert.Equal("webserver.local", cfg.ServerHost)
	assert.Equal(uint16(8080), cfg.ServerPort)
	assert.Equal("/var/www/site1", cfg.ServerFilePath)

	assert.Equal("/var/migrations", cfg.MigrationsPath)
}
