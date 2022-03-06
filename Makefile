
.PHONY: help
help: ## Show this help
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## Run the golangci linter
	@docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint golangci-lint run

.PHONY: serve
serve: ## Run the server
	@docker run --rm -v $$(pwd):/app -w /app -p 8080:8080 golang:1.17.8-buster go run main.go
