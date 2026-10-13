.DEFAULT_GOAL := help

EXEC = docker compose exec app
ARTISAN = $(EXEC) go run . artisan

.PHONY: help up down build shell logs run artisan migrate seed fresh test fmt vet lint

help: ## List available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

up: ## Start all containers in the background
	@docker compose up -d

down: ## Stop and remove containers
	@docker compose down

build: ## Rebuild the app image
	@docker compose build

shell: ## Open a shell in the app container
	@$(EXEC) bash

logs: ## Tail app container logs
	@docker compose logs -f app

run: ## Restart the app container (air rebuilds the binary on start)
	@docker compose restart app

artisan: ## Run an artisan command, e.g. make artisan cmd="route:list"
	@$(ARTISAN) $(cmd)

migrate: ## Run database migrations
	@$(ARTISAN) migrate

seed: ## Run database seeders
	@$(ARTISAN) db:seed

fresh: ## Drop all tables, re-run migrations and seed
	@$(ARTISAN) migrate:fresh --seed

test: ## Run the test suite
	@$(EXEC) go test ./...

fmt: ## Format all Go files with gofmt
	@$(EXEC) gofmt -w .

vet: ## Run go vet static analysis
	@$(EXEC) go vet ./...

lint: ## Fail if any file needs gofmt, then run go vet
	@files=$$($(EXEC) gofmt -l .); \
	if [ -n "$$files" ]; then echo "needs gofmt (run make fmt):"; echo "$$files"; exit 1; fi
	@$(EXEC) go vet ./...
	@echo "lint passed"
