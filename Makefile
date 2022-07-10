GOBIN ?= $(shell pwd)/bin
PATH := $(GOBIN):$(PATH)
export

GOFUMPT = $(GOBIN)/gofumpt
GOLANGCI_LINT = $(GOBIN)/golangci-lint
GOIMPORTS = $(GOBIN)/goimports-reviser

$(GOFUMPT):
	cd tools && go install mvdan.cc/gofumpt

$(GOIMPORTS):
	cd tools && go install github.com/incu6us/goimports-reviser

$(GOLANGCI_LINT):
	cd tools && go install github.com/golangci/golangci-lint/cmd/golangci-lint


# Directories containing independent Go modules.
MODULE_DIRS = ./backend_service ./mongo_service ./postgres_service

# Many Go tools take file globs or directories as arguments instead of packages.
GO_FILES := $(shell \
	find . '(' -path '*/.*' -o -path './vendor' ')' -prune \
	-o -name '*.go' -print | cut -b3-)

.PHONY: gen fmt lint tidy test dev clean help

default: help

fmt: $(GOFUMPT) $(GOIMPORTS) ## Format source files.
	@echo "formatting via gofumpt and goimports-reviser"
	@$(foreach file,$(GO_FILES),(echo "fmt $(file)" && $(GOFUMPT) -e -w $(file) && $(GOIMPORTS) -project-name github.com/kucera-lukas/micro-backends -file-path $(file)) &&) true

lint: $(GOLANGCI_LINT)  ## Lint source files.
	@echo "linting via golangci-lint"
	@$(GOLANGCI_LINT) run --config ./.golangci-lint.yml ./...

tidy:  ## Tidy module dependencies.
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && echo "tidy $(dir)" && go mod tidy -compat=1.18) &&) true

test:  ## Run module tests.
	@$(foreach dir,$(MODULE_DIRS),(cd $(dir) && echo "test $(dir)" && go test -race ./...) &&) true

gen-backend:  ## Generate files for backend_service.
	@echo "generating gRPC files for backend_service"
	@cd backend_service/proto && protoc backend.proto --go_out=. --go-grpc_out=.

gen-mongo:  ## Generate files for mongo_service.
	@echo "generating gRPC files for mongo_service"
	@echo $(PATH)
	@cd mongo_service/proto && protoc mongo.proto --go_out=. --go-grpc_out=.

gen-postgres:  ## Generate files for postgres_service.
	@echo "generating gRPC files for postgres_service"
	@cd postgres_service/proto && protoc postgres.proto --go_out=. --go-grpc_out=.

gen:  ## Generate files for all services.
	@make gen-backend
	@make gen-mongo
	@make gen-postgres

dev:  ## Start development server.
	@echo "starting docker compose network"
	@docker-compose up -d

dev-build:  ## Build and start development server.
	@echo "building docker compose network"
	@docker-compose up --build -d

clean:  ## Cleanup binary files.
	@echo "cleaning up 'bin'"
	@rm -rf bin/*

help: ## Display help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
