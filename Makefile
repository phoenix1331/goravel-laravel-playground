.DEFAULT_GOAL := help

EXEC = docker compose exec app
ARTISAN = $(EXEC) go run . artisan

.PHONY: help up down build shell logs run artisan migrate seed fresh test test-mysql fmt vet lint docs-dev docs-build

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

test: ## Run the test suite against in-memory SQLite (.env.testing)
	@$(EXEC) go test ./... -count=1 --env=.env.testing

test-mysql: ## Run the test suite production-like: MySQL (goravel_test) and the Redis cache (db 1)
	@docker compose exec -e DB_CONNECTION=mysql -e DB_HOST=mysql -e DB_PORT=3306 -e DB_DATABASE=goravel_test -e DB_USERNAME=goravel -e DB_PASSWORD=secret -e CACHE_STORE=redis -e REDIS_HOST=redis -e REDIS_DB=1 app go test ./... -count=1 -p 1 --env=.env.testing

fmt: ## Format all Go files with gofmt
	@$(EXEC) gofmt -w .

vet: ## Run go vet static analysis
	@$(EXEC) go vet ./...

lint: ## Fail if any file needs gofmt, then run go vet
	@files=$$($(EXEC) gofmt -l .); \
	if [ -n "$$files" ]; then echo "needs gofmt (run make fmt):"; echo "$$files"; exit 1; fi
	@$(EXEC) go vet ./...
	@echo "lint passed"

docs-dev: ## Serve the VitePress docs with hot reload (host Node, http://localhost:5173)
	@npm run docs:dev

docs-build: ## Build the static docs site into docs/.vitepress/dist
	@npm run docs:build
