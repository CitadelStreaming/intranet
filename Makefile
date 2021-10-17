.PHONY: all covtest fmt

export CGO_ENABLED=0
export GOARCH=amd64
TAGS:="-integration"

all: test vet citadel_intranet

citadel_intranet:
	go build -o citadel_intranet src/main.go

deployment_executable: citadel_intranet
	strip citadel_intranet

fmt:
	go fmt ./...

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

intgtestreport: intgtest
	go tool cover -html=coverage.out

mock:
	./bin/mocks.sh

docker_image: deployment_executable
	./bin/docker.sh

clean:
	rm -rf coverage.out citadel_intranet gomock_reflect_*
	find . -name mock -a -type d -exec rm -rf {} \;
