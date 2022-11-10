BINARY_NAME=main
 
build:
		go mod tidy
		go build -o ${BINARY_NAME} main.go

run:
	  go run main.go

runWithJSON:
		go run main.go -j
 
clean:
		go clean
		# rm -f ${BINARY_NAME}
