prepare:
	go mod tidy

run:
	go run -race main.go handler.go github.go

build+linux:
	mkdir -p bin/linux
	GOOS=linux GOARCH=amd64 go build  -o bin/linux/private-ghp main.go handler.go github.go

build+docker: build+linux
	docker build -t private-ghp:latest .