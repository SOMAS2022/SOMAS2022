BINARY_NAME=cmd/main.out
SOURCE_DIR=pkg/infra

PACKAGES=$(shell go list ./pkg/infra | grep -v 'tests')

all: run

build:
		cd ${SOURCE_DIR}; go mod tidy
		go build -o ${BINARY_NAME} infra

run: build
		${BINARY_NAME}

runWithID: build
		${BINARY_NAME} -i ${ID}

runWithJSON: build
		${BINARY_NAME} -j

runDebug: build
		${BINARY_NAME} -d
clean:
		go clean
		rm -r logs/*
		rm -rf ${BINARY_NAME}

# formatting and linting
fmt:
	gofmt -s -w .
# change to `run ./pkg/*` after agents are implemented
# should just be `run`, but seems to be problems with go.work  
# nb: using cd to /pkg/infra dosen't fix this (on wsl2)
check:
	golangci-lint -v run ./pkg/infra

# for future testing
unit_test:
	go test $(PACKAGES)

test:
	go test ./pkg/infra -covermode=atomic

test_race:
	go test ./pkg/infra --race
