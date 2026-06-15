APP         := ims
CMD         := .
BIN         := ./$(APP)
SWAG        := swag

.PHONY: all build run test clean tidy lint swagger db-start db-stop migrate help

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

## db-start: start PostgreSQL container via Docker Compose
db-start:
	docker-compose up -d postgres
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker exec ims-postgres pg_isready -U $(DB_USER) -d ims_db > /dev/null 2>&1; do sleep 1; done
	@echo "PostgreSQL is ready."

## db-stop: stop and remove PostgreSQL container
db-stop:
	docker-compose down

## migrate: run Liquibase migrations via Docker Compose
migrate:
	docker-compose --profile migrate run --rm liquibase

## help: list available targets
help:
	@grep -E '^##' Makefile | sed 's/## /  /'
