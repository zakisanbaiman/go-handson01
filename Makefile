.PHONY: help build build-local up down logs ps test
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

help: ## Show help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'