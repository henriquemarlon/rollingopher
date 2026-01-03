-include .env

START_LOG = @echo "======================= START OF LOG ======================="
END_LOG   = @echo "======================== END OF LOG ========================"

.PHONY: build-echo
build-echo: ## Build the echo example
	$(START_LOG)
	@cartesi build -c cartesi.echo.toml
	$(END_LOG)

.PHONY: build-handling-assets
build-handling-assets: ## Build the handling-assets example
	$(START_LOG)
	@cartesi build -c cartesi.handling-assets.toml
	$(END_LOG)

.PHONY: test-echo
test-echo: build-echo ## Run the echo tests
	$(START_LOG)
	@pnpm install
	@pnpm vitest run examples/echo/echo.test.ts
	$(END_LOG)

.PHONY: test-handling-assets
test-handling-assets: build-handling-assets ## Run the handling-assets tests
	$(START_LOG)
	@pnpm install
	@pnpm vitest run examples/handling-assets/handling-assets.test.ts
	$(END_LOG)

.PHONY: fmt
fmt: ## Format all code (Go + TypeScript)
	$(START_LOG)
	@gofmt -w .
	@pnpm exec prettier --log-level silent --write "**/*.ts"
	@echo "Formatting completed"
	$(END_LOG)

.PHONY: fmt-go
fmt-go: ## Format Go code only
	$(START_LOG)
	@gofmt -w .
	@echo "Go formatting completed"
	$(END_LOG)

.PHONY: fmt-ts
fmt-ts: ## Format TypeScript code only
	$(START_LOG)
	@pnpm exec prettier --log-level silent --write "**/*.ts"
	@echo "TypeScript formatting completed"
	$(END_LOG)

.PHONY: clean
clean: ## Clean all build artifacts
	$(START_LOG)
	@rm -rf .cartesi
	@rm -rf build
	@echo "Clean completed"
	$(END_LOG)

.PHONY: help
help: ## Show help for each of the Makefile recipes
	@echo "Available commands:"
	@awk '/^[a-zA-Z0-9_-]+:.*##/ { \
		split($$0, parts, "##"); \
		split(parts[1], target, ":"); \
		printf "  \033[36m%-30s\033[0m %s\n", target[1], parts[2] \
	}' $(MAKEFILE_LIST)
