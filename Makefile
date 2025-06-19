
SRC := main.go
BIN := trans
BUILD-DIR := build

.PHONY: build run clean release

build: build-linux-amd64 build-darwin-amd64 build-darwin-arm64 build-windows-amd64

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD-DIR)/$(BIN)-linux-amd64/$(BIN) main.go

build-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD-DIR)/$(BIN)-darwin-amd64/$(BIN) main.go

build-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD-DIR)/$(BIN)-darwin-arm64/$(BIN) main.go

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD-DIR)/$(BIN)-windows-amd64/$(BIN).exe main.go

release: build
	sh ./release.sh

run:
	go run main.go

clean:
	rm -r $(BUILD-DIR)
	rm -r release