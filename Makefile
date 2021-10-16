.PHONY: all covtest

TAGS:="config,-integration"

all: test vet
	go build src/main.go

vet:
	go vet ./...

test: mock
	go test --tags="${TAGS}" ./...

covtest: mock
	go test -coverprofile=coverage.out --tags="${TAGS}" ./...

covreport: covtest
	go tool cover -html=coverage.out

intgtest: mock
	./bin/intgtest.sh

mock:
	./bin/mocks.sh

clean:
	rm -rf coverage.out main gomock_reflect_*
	find . -name mock -a -type d -exec rm -rf {} \;
