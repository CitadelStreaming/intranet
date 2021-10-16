// +build integration

package db_test

import (
	"os"
	"testing"

	"citadel_intranet/src/config"
	"citadel_intranet/src/db"

	"github.com/stretchr/testify/assert"
)

func TestRunMigrations(t *testing.T) {
	assert := assert.New(t)

	callback := func() {
		wd, err := os.Getwd()
		assert.Nil(err)

		cfg := config.Config{
			DbHost: "localhost",
			DbPort: 3306,
			DbUser: "root",
			DbPass: "pass",
			DbName: "testbed",

			MigrationsPath: wd + "/../../migrations/",
		}

		db := db.NewDatabaseClient(cfg)
		assert.NotNil(db)
		defer db.Close()
		db.Migrate(cfg.MigrationsPath)
	}

	callback()

	// We shouldn't have any new migrations to run since we ran them once
	// already, but this should exercise a slightly different code path.
	callback()
}

func TestBadConnection(t *testing.T) {
	assert := assert.New(t)
	(func() {
		defer (func() {
			assert.NotNil(recover())
		})()

		wd, err := os.Getwd()
		assert.Nil(err)

		cfg := config.Config{
			DbHost: "localhost",
			DbPort: 3306,
			DbUser: "root",
			DbPass: "notTheCorrectPassword",
			DbName: "testbed",

			MigrationsPath: wd + "/../../migrations/",
		}

		// This should be a fatal error due to not being able to connect
		client := db.NewDatabaseClient(cfg)
		client.Migrate(cfg.MigrationsPath)
		defer client.Close()
		assert.False(true, "The code didn't hit a fatal error")
	})()
}
