BINARY_NAME=cmd/main.out
SOURCE_DIR=pkg/infra

${BINARY_NAME}:
		cd ${SOURCE_DIR}; go mod tidy
		go build -o ${BINARY_NAME} infra

build: ${BINARY_NAME}

run: build
	  ${BINARY_NAME}

runWithJSON: build
		${BINARY_NAME} -j
 
clean:
		go clean
		rm -rf ${BINARY_NAME}
