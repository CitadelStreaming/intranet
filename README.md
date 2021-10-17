# 🎜 Citadel Streaming Intranet

The intranet site for Citadel Streaming, meant to provide an easy interface to
start putting together new albums, collaborate on naming and track thoughts, and
see everything come together before publishing the album.

## Environment Variables
 
* `SECRET_LOCATION` **Unused**
* `DB_HOST` The hostname or IP address for the database.
* `DB_PORT` Port number for the database.
* `DB_NAME` The name of the database to use.
* `DB_USER` Username for the database.
* `DB_PASS` Password for the database.
* `SERVER_HOST` The hostname to use for the server (where we should bind to).
* `SERVER_PORT` The port number to serve on.
* `SERVER_PATH` The location to serve static files from.
* `MIGRATIONS` The location to look for database migrations in.

## Building, testing, and more

Everything is currently done via `make`, each of the targets are there to make
your life easier, and allow for quick iterations. Available targets (at the time
or writing) are:

### `all`

Run through tests, vet, and finally build the intranet application.

### `fmt`

Run the go formatter.

### `test`

Run unit tests.

### `covtest`

Run unit tests with coverage.

### `covreport`

Runs unit tests with coverage and generates an HTML report.

### `intgtest`

Run unit _and_ integration tests. Replies on setting up a docker MariaDB
database.

### `intgtestreport`

Runs the integration tests and generates an HTML report of the coverage numbers.

### `mock`

Generates interface mocks automatically. All mocks will be placed in a `mock`
package under their current directory structure. These files should **NOT** be
checked in!

### `vet`

Run through the linter.

### `clean`

Clean up all mocks, coverage reports, and the main application.

## Migrations

Migrations are responsible for putting (most) of the database into an expected
state. The migrator sets up the migrations table, which it relies explicitly on,
and as such is desirable to have around _before_ any migrations are run.

Migrations files are all located in `migrations/` and should all be MariaDB
syntax. They will be run in order, and as such should be named in the form of
`YYYYMMDD_XXX.sql` with the date and then the number (incrementing) of the
migration for that day. Migrations are considered a zeroith-order citizen,
meaning that if they cannot be completed for any reason, the application _will_
panic.

Some things to note:  
* Migrations MUST be added in order. You _cannot_ come back and add a migration
  that takes place before the most recent migration. That will cause a panic.
* Migrations MUST not be changed after they are migrated. We checksum (sha1)
  each migration _before_ running it and store that hash. If the hash of the
  file _ever_ changes after we have run the migration, we panic.
* We allow multiple queries in a single file, however the MySQL driver for go
  has issues with that, for that reason we split each separate query apart and
  run them individually. Don't worry, however, as all migrations are run within
  a transaction. If we return without a panic for a migration, the migration has
  been completed and recorded without error.
* A SQL error encountered in a migration will cause a panic.
