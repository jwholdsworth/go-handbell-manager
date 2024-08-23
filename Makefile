run: build
	sudo ./go-handbell-manager

build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" ./...
