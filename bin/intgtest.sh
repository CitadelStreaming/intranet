#!/usr/bin/env bash

./bin/oci-exec run --detach -p 3306:3306 --rm --env=MYSQL_ROOT_PASSWORD=pass mariadb:latest
container=$(./bin/oci-exec ps -q | tail -1)

# Wait for the database to be up and running. Apparently there isn't anything
# like `pg_isready` that we can use here...
for i in {0..10}; do
    mysql --user=root --password=pass --port=3306 --host=127.0.0.1 -e "CREATE DATABASE IF NOT EXISTS testbed" || sleep 1 && continue
    break
done

mysql --user=root --password=pass --port=3306 --host=127.0.0.1 -e "SELECT VERSION()" testbed && \
go test -coverprofile=coverage.out --tags="integration" ./...

./bin/oci-exec kill ${container}
