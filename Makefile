.PHONY: all covtest

all: test vet
	go build src/main.go

vet:
	go vet ./...


test:
	go test ./...

covtest:
	go test -cover ./...
