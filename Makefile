.PHONY: run build test

run: build
	@./bin/api

build:
	@mkdir -p bin
	@go build -o bin/api ./cmd/app

test:
	@go test -v ./...
