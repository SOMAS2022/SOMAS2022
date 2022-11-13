BINARY_NAME?=cmd/main.out
SOURCE_DIR=pkg/infra
TEAM?=default
LDFLAGS=-ldflags="-X 'infra/game/stages.mode=${TEAM}'"

all: run

build:
		cd ${SOURCE_DIR}; go mod tidy
		go build $(LDFLAGS) -o ${BINARY_NAME} infra 

run: build
		${BINARY_NAME}

runWithJSON: build
		${BINARY_NAME} -j
 
clean:
		go clean
		rm -rf ${BINARY_NAME}
