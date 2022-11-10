BINARY_NAME=cmd/main.out
SOURCE_DIR=pkg/infra
ROOT_DIR=..

${BINARY_NAME}:
		cd ${ROOT_DIR}/${SOURCE_DIR}; go mod tidy
		go build -o ${ROOT_DIR}/${BINARY_NAME} infra

build: ${BINARY_NAME}

run: ${BINARY_NAME}
	  cd ${ROOT_DIR}; ${BINARY_NAME}

runWithJSON: ${BINARY_NAME}
		cd ${ROOT_DIR}; ${BINARY_NAME} -j
 
clean:
		cd ${BUILD_DIR}; go clean
		rm -rf ${ROOT_DIR}/${BINARY_NAME}
