GOLANGCI_LINT_VERSION := v2.13.2

.PHONY: help run build test lint lint-install tidy compose-up compose-down compose-restart compose-logs migrate-up migrate-down

help:
	@echo "make run              - run the api locally"
	@echo "make build            - build the api binary into bin/"
	@echo "make test             - run unit tests"
	@echo "make lint             - run golangci-lint"
	@echo "make lint-install     - install golangci-lint $(GOLANGCI_LINT_VERSION)"
	@echo "make tidy             - tidy go modules"
	@echo "make compose-up       - start postgres and the app"
	@echo "make compose-down     - stop and remove the containers"
	@echo "make compose-restart  - rebuild and restart the containers"
	@echo "make compose-logs     - follow the container logs"
	@echo "make migrate-up       - apply the database migration"
	@echo "make migrate-down     - drop the database migration"

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api
	go build -o bin/migrate ./cmd/migrate

test:
	go test ./... -v -cover

lint:
	golangci-lint run ./...

lint-install:
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

tidy:
	go mod tidy

compose-up:
	docker compose up -d

compose-down:
	docker compose down

compose-restart:
	docker compose up -d --build

compose-logs:
	docker compose logs -f

migrate-up:
	go run ./cmd/migrate -m up

migrate-down:
	go run ./cmd/migrate -m down
