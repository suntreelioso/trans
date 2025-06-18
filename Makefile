
.PHONY: build run clean

build:
	go build -o trans main.go

run:
	go run main.go

clean:
	rm trans
