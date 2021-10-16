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

	assert.Equal("/var/migrations", cfg.MigrationsPath)
}

func TestLoadConfigSetValues(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(os.Setenv(config.ENV_DATABASE_HOST, "database.local"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_USER, "bobby"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_PASS, "tables"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_PORT, "23306"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_NAME, "db1"))

	assert.Nil(os.Setenv(config.ENV_MIGRATIONS_PATH, "/opt/citadel/migrations"))

	cfg := config.LoadConfig()

	assert.Equal("database.local", cfg.DbHost)
	assert.Equal("bobby", cfg.DbUser)
	assert.Equal("tables", cfg.DbPass)
	assert.Equal(uint16(23306), cfg.DbPort)
	assert.Equal("db1", cfg.DbName)

	assert.Equal("/opt/citadel/migrations", cfg.MigrationsPath)
}

func TestLoadConfigSetValuesInvalidPort(t *testing.T) {
	assert := assert.New(t)

	assert.Nil(os.Setenv(config.ENV_DATABASE_HOST, "database.local"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_USER, "bobby"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_PASS, "tables"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_PORT, "NotANumber"))
	assert.Nil(os.Setenv(config.ENV_DATABASE_NAME, "db1"))

	assert.Nil(os.Setenv(config.ENV_MIGRATIONS_PATH, "/var/migrations"))

	cfg := config.LoadConfig()

	assert.Equal("database.local", cfg.DbHost)
	assert.Equal("bobby", cfg.DbUser)
	assert.Equal("tables", cfg.DbPass)
	assert.Equal(uint16(3306), cfg.DbPort)
	assert.Equal("db1", cfg.DbName)

	assert.Equal("/var/migrations", cfg.MigrationsPath)
}
