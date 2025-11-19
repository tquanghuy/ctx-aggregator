.PHONY: all test lint fmt clean

all: test lint

test:
	go test -v -race ./...

lint:
	golangci-lint run ./...

fmt:
	go fmt ./...

clean:
	go clean
