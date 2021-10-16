# 🎜 Citadel Streaming Intranet

The intranet site for Citadel Streaming, meant to provide an easy interface to
start putting together new albums, collaborate on naming and track thoughts, and
see everything come together before publishing the album.

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
