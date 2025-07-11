buildbin:
	GOOS=linux GOARCH=amd64 go build -o ./build/mayohttp .
run:
	go run .
watch:
	gow run .
