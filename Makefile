.PHONY: help build build-local up down logs ps test lint ci-local
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build: ## Build docker image to deploy
	docker build -t zakisanbaiman/gotodo:${DOCKER_TAG} --target deploy .

build-local: ## Build docker image for local development
	docker compose build --no-cache

up: ## Run docker container
	docker compose up -d

down: ## Stop docker container
	docker compose down

logs: ## Show logs
	docker compose logs -f

ps: ## Show docker ps
	docker compose ps

test: ## Run test
	go test -race -shuffle=on ./...

migrate: ## Run migrate
	mysqldef -u todo -p todo -h 127.0.0.1 -P 33306 todo < _tools/mysql/schema.sql

dry-migrate: ## Run dry migrate
	mysqldef -u todo -p todo -h 127.0.0.1 -P 33306 todo < _tools/mysql/schema.sql --dry-run

generate:
	go generate ./...

lint: ## Run golangci-lint locally
	golangci-lint run ./...

ci-local: ## Run GitHub Actions locally using act
	DOCKER_HOST=unix:///Users/n_suzaki/.orbstack/run/docker.sock act -W .github/workflows/golangci.yml --container-architecture linux/amd64

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'