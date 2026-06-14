APP     := ims
CMD     := .
BIN     := ./$(APP)
SWAG    := swag

.PHONY: all build run test clean tidy lint swagger help

all: tidy swagger build

## build: compile the binary
build:
	go build -o $(BIN) $(CMD)

## run: build and start the server on :8080
run: build
	$(BIN)

## test: run all unit tests with race detector
test:
	go test -race -v ./...

## tidy: download and verify dependencies
tidy:
	go mod tidy

## swagger: regenerate Swagger docs from annotations
swagger:
	$(SWAG) init -g main.go -o docs

## lint: run go vet across all packages
lint:
	go vet ./...

## clean: remove the compiled binary
clean:
	rm -f $(BIN)

## help: list available targets
help:
	@grep -E '^##' Makefile | sed 's/## /  /'
