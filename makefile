BINARY_NAME=ams
WORK_DIR=assets/music/
CRITERIA=artist

run: build
	./${BINARY_NAME}

e2e1: build
	./${BINARY_NAME} ${WORK_DIR} ${CRITERIA}

e2e2 : WORK_DIR=assets/music/more/
e2e2 : build
	./${BINARY_NAME} ${WORK_DIR} ${CRITERIA}

build:
	@go build -o ${BINARY_NAME} cmd/${BINARY_NAME}/main.go

clean:
	@go clean
	@rm ${BINARY_NAME}
