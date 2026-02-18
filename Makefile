.PHONY: build run test lint clean

build:
	go build -o bin/worker ./cmd/worker

run:
	go run ./cmd/worker

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
