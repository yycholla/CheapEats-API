.PHONY: help
help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: install
install: ## Install dependencies
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest

.PHONY: swagger
swagger: ## Generate swagger documentation
	swag init -g cmd/api/main.go --output docs

.PHONY: build
build: swagger ## Build the application
	go build -o bin/api cmd/api/main.go

.PHONY: run
run: swagger ## Run the application locally
	go run cmd/api/main.go

.PHONY: docker-build
docker-build: ## Build docker image
	docker build -t cheapeats-api .

.PHONY: docker-up
docker-up: ## Start docker containers
	docker-compose up -d

.PHONY: docker-down
docker-down: ## Stop docker containers
	docker-compose down

.PHONY: docker-logs
docker-logs: ## View docker logs
	docker-compose logs -f

.PHONY: test
test: ## Run tests
	go test -v ./...

.PHONY: lint
lint: ## Run linter
	golangci-lint run

.PHONY: clean
clean: ## Clean build artifacts
	rm -rf bin/ docs/

.PHONY: db-migrate
db-migrate: ## Run database migrations
	go run cmd/api/main.go migrate

.PHONY: dev
dev: swagger docker-up ## Start development environment
	air