CONTAINER_RUNTIME=$(shell which docker)

prepare:
	go mod tidy

run:
	go run -race main.go handler.go github.go

build:
	mkdir -p bin/generic
	go build -o bin/generic/private-ghp main.go handler.go github.go

build+darwin:
	mkdir -p bin/darwin
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/private-ghp main.go handler.go github.go

build+linux:
	mkdir -p bin/linux
	GOOS=linux GOARCH=amd64 go build -o bin/linux/private-ghp main.go handler.go github.go

build+docker:
	${CONTAINER_RUNTIME} build -t private-ghp:latest .