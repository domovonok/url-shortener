.DEFAULT_GOAL := run

MAIN           ?= ./cmd/app/main.go
BIN            ?= app
MIGRATIONS_DIR ?= ./migrations

GOOSE_CMD = goose -dir $(MIGRATIONS_DIR) $(GOOSE_DRIVER) $(GOOSE_DBSTRING)

.PHONY: up down migrate migrate-down run build

up:
	docker compose up -d

down:
	docker compose down

migrate:
	$(GOOSE_CMD) up

migrate-down:
	$(GOOSE_CMD) down

run: up migrate
	go run $(MAIN)

build:
	go build -o $(BIN) $(MAIN)
