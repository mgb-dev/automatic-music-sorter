BINARY_NAME=ams

run: build
	./${BINARY_NAME}

build:
	@go build -o ./${BINARY_NAME} cmd/${BINARY_NAME}/main.go

test:
	go test -v

clean:
	@go clean
	@rm ./${BINARY_NAME}
