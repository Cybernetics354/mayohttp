buildbin:
	GOOS=linux GOARCH=amd64 go build -o ./build/mayohttp .
run:
	go run cmd/mayohttp/main.go
watch:
	gow run cmd/mayohttp/main.go
lint:
	golangci-lint run
