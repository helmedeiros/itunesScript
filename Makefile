# Quality gates for the Apple Music terminal player.
# Every target must stay green on every commit.

GO        ?= go
GOTESTSUM ?= gotestsum
GOLANGCI  ?= golangci-lint

.DEFAULT_GOAL := ci

.PHONY: ci
ci: fmt-check vet lint test ## Full gate: run in CI and before every commit

.PHONY: fmt
fmt: ## Format the code in place
	$(GOLANGCI) fmt

.PHONY: fmt-check
fmt-check: ## Fail if any file is not formatted
	$(GOLANGCI) fmt --diff

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	$(GOLANGCI) run

.PHONY: test
test: ## Run unit tests with the race detector
	$(GOTESTSUM) --format dots-v2 -- -race -count=1 ./...

.PHONY: cover
cover: ## Report test coverage
	$(GO) test -race -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out | tail -1

.PHONY: integration
integration: ## Run osascript integration tests (requires macOS + Music.app)
	$(GOTESTSUM) --format testname -- -race -tags=integration -run Integration ./...

.PHONY: build
build: ## Build all binaries into ./bin
	$(GO) build -o bin/ ./cmd/...

.PHONY: tidy
tidy: ## Tidy module dependencies
	$(GO) mod tidy

.PHONY: help
help: ## List available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
