.PHONY: build
build: ## Build a version
	make swag
	go build -v ./cmd/goblog

.PHONY: test
test: ## Run tests
	go test -v

.PHONY: clean
clean: ## Remove temporary files
	go clean

.PHONY: dev
dev: ## Go Run
	export DATABASE_DSN="host=localhost user=postgres password=postgres dbname=app port=5432 sslmode=disable" && go run cmd/goblog/main.go

.PHONY: swag
swag: ## Update swagger.json
	swag init -g ./cmd/goblog/main.go

.PHONY: swag-fmt
swag-fmt: ## Formatter for GoDoc (Swagger)
	swag fmt -g ./cmd/goblog/main.go

.PHONY: docker-up
docker-up: ## Start Docker-Compose Container with app & database
	docker-compose -f build/docker-compose.yml up -d --build

.PHONY: docker-down
docker-down: ## Down Docker-Compose Containers
	docker-compose -f build/docker-compose.yml down

.PHONY: docker-database-up
docker-database-up: ## Start Docker-compose Container with only database service
	docker-compose -f build/docker-compose.yml up database -d

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build