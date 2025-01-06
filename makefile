BINARY_NAME=ams
SYMLINK_DESTINATION=${HOME}/.local/bin/${BINARY_NAME}

run: build
	./${BINARY_NAME}

build:
	@go build -o ./${BINARY_NAME} cmd/${BINARY_NAME}/main.go

test:
	go test -v

install: build
	@ln -s -f ${PWD}/ams ${SYMLINK_DESTINATION}

clean:
	@go clean
	@rm ./${BINARY_NAME}
