.PHONY: all covtest

all: test vet
	go build src/main.go

vet:
	go vet ./...


test: mock
	go test ./...

covtest: mock
	go test -coverprofile=coverage.out ./...

covreport: covtest
	go tool cover -html=coverage.out

mock:
	./bin/mocks.sh
