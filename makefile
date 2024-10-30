ARTIFACT_NAME=ams

run:
	@go run cmd/${ARTIFACT_NAME}/main.go
build:
	@go build -o ${BINARY_NAME} cmd/${BINARY_NAME}/main.go
