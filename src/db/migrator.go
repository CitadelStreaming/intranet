package db

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func migrate(db *sql.DB, migrationsPath string) {
	ensureMigrationsTableExists(db)

	items, err := os.ReadDir(migrationsPath)

	if err != nil {
		logrus.Panic("Unable to run migrations for given path ", migrationsPath, " ", err.Error())
	}

	rows := retrieveCompletedMigrations(db)
	rowsClosed := false

	defer (func() {
		if !rowsClosed {
			rows.Close()
		}
	})()

	var migrationName string
	var migrationSha string
	for _, item := range items {
		// Skip dotfiles
		if item.Name()[0] == '.' {
			continue
		}

		if rowsClosed {
			migrationBody, fileHash := fileSha(migrationsPath, item.Name())
			executeMigration(db, string(migrationBody), fileHash, item.Name())
		} else if rows.Next() {

			rows.Scan(&migrationName, &migrationSha)
			if migrationName != item.Name() {
				logrus.Panic("Unexpected migration found ", item.Name(), " expecting ", migrationName)
			}

			_, fileHash := fileSha(migrationsPath, item.Name())
			if migrationSha != fileHash {
				logrus.Panic("Migration has been modified since it was applied! ", migrationName, " ", migrationSha)
			}

			// We don't have a migration to run here.
		} else {
			rowsClosed = true
			migrationBody, fileHash := fileSha(migrationsPath, item.Name())
			executeMigration(db, string(migrationBody), fileHash, item.Name())
		}
	}
}

func ensureMigrationsTableExists(db *sql.DB) {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS migrations(
            id BIGINT PRIMARY KEY NOT NULL AUTO_INCREMENT,
            name VARCHAR(255) UNIQUE NOT NULL DEFAULT '',
            checksum VARCHAR(40) NOT NULL DEFAULT ''
        )
    `)

	if err != nil {
		logrus.Panic("Unable to create migrations table ", err.Error())
	}
}

func retrieveCompletedMigrations(db *sql.DB) *sql.Rows {
	rows, err := db.Query(`
        SELECT
            name,
            checksum
        FROM migrations
        ORDER BY
            name ASC
    `)
	if err != nil {
		logrus.Panic("Failed to load migrations: ", err.Error())
	}

	return rows
}

func fileSha(dir string, file string) ([]byte, string) {
	migration := dir + "/" + file
	f, err := os.OpenFile(migration, os.O_RDONLY, 0755)

	if err != nil {
		logrus.Panic("Unable to open file ", migration)
	}
	bytes, err := ioutil.ReadAll(f)
	f.Close()

	return bytes, fmt.Sprintf("%x", sha1.Sum(bytes))
}

func executeMigration(db *sql.DB, migrationQuery string, hash string, migrationName string) {
	logrus.WithField("query", migrationQuery).Info("Running migration: ", migrationName)

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		logrus.Panic("Migration error: ", err.Error())
	}

	// MySQL driver doesn't support multiple queries in a single Exec, so we
	// split them apart and run each individually.
	queries := strings.Split(migrationQuery, ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		logrus.Info(query)
		_, err = tx.Exec(query)
		if err != nil {
			tx.Rollback()
			logrus.Panic("Migration error: ", err.Error())
		}
	}

	logrus.Info("Writing checksum: ", hash, " for migration: ", migrationName)
	_, err = tx.Exec(`
        INSERT INTO migrations(
            name,
            checksum
        )
        VALUES(
            ?,
            ?
        )
    `,
		migrationName,
		hash)

	if err != nil {
		tx.Rollback()
		logrus.Panic("Migration recording error: ", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		logrus.Panic("Transaction commit error: ", err.Error())
	}

}
