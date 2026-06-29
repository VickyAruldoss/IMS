APP         := ims
CMD         := .
BIN         := ./$(APP)
SWAG        := swag

# Load environment variables from .env (if present) and export them to recipes
-include .env
export

# Fallback defaults if .env is missing
DB_USER     ?= ims_user
DB_NAME     ?= ims_db

PGADMIN_PORT ?= 5050

.PHONY: all build run test clean tidy lint swagger db-start db-stop migrate pgadmin pgadmin-stop help

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
	@until docker exec ims-postgres pg_isready -U $(DB_USER) -d $(DB_NAME) > /dev/null 2>&1; do sleep 1; done
	@echo "PostgreSQL is ready."

## db-stop: stop and remove PostgreSQL container
db-stop:
	docker-compose down

## migrate: run Liquibase migrations via Docker Compose
migrate:
	docker-compose --profile migrate run --rm liquibase

## pgadmin: start pgAdmin4 container via Docker Compose
pgadmin:
	docker-compose --profile pgadmin up -d pgadmin
	@echo "pgAdmin4 available at http://localhost:$(PGADMIN_PORT)"

## pgadmin-stop: stop and remove pgAdmin4 container
pgadmin-stop:
	docker-compose --profile pgadmin rm -sf pgadmin

## help: list available targets
help:
	@grep -E '^##' Makefile | sed 's/## /  /'
