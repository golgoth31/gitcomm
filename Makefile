.PHONY: build test lint format clean install

build:
	go build -o gitcomm ./cmd/gitcomm

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint:
	golangci-lint run

format:
	go fmt ./...
	goimports -w .

clean:
	rm -f gitcomm coverage.out

install:
	go install ./cmd/gitcomm
