GOBIN ?= $(shell pwd)/bin
PATH := $(GOBIN):$(PATH)
GOFUMPT_SPLIT_LONG_LINES := on
export

GQLGEN = $(GOBIN)/gqlgen
GOFUMPT = $(GOBIN)/gofumpt
GOLANGCI_LINT = $(GOBIN)/golangci-lint
GOIMPORTS = $(GOBIN)/goimports-reviser
GOLINES = $(GOBIN)/golines

$(GQLGEN):
	cd tools && go install github.com/99designs/gqlgen

$(GOFUMPT):
	cd tools && go install mvdan.cc/gofumpt

$(GOIMPORTS):
	cd tools && go install github.com/incu6us/goimports-reviser

$(GOLANGCI_LINT):
	cd tools && go install github.com/golangci/golangci-lint/cmd/golangci-lint

$(GOLINES):
	cd tools && go install github.com/segmentio/golines


# Directories containing independent Go modules.
MODULE_DIRS = ./backend_service ./mongo_service ./postgres_service ./tools

# Many Go tools take file globs or directories as arguments instead of packages.
GO_FILES := $(shell \
	find . '(' -path '*/.*' -o -path './vendor' ')' -prune \
	-o -name '*.go' -print | cut -b3-)

.PHONY: fmt lint tidy test gen-backend gen-mongo gen-postgres gen build dev dev-build clean help

default: help

fmt: $(GOLINES) $(GOFUMPT) $(GOIMPORTS)  ## Format source files.
	@echo "formatting via golines"
	@$(GOLINES) --max-len=80 -w .
	@echo "formatting via gofumpt"
	@$(GOFUMPT) -e -w -extra .
	@echo "formating via goimports-reviser"
	@$(foreach file,$(GO_FILES),( \
		$(GOIMPORTS) \
			-project-name github.com/kucera-lukas/micro-backends \
			-file-path $(file) \
		  ) && \
		) true

lint: $(GOLANGCI_LINT)  ## Lint source files.
	@echo "linting via golangci-lint"
	@$(foreach dir,$(MODULE_DIRS),( \
		cd $(dir) && \
		echo "lint $(dir)" && \
		$(GOLANGCI_LINT) run --config ../.golang-ci.lint.yml ./...) && \
		) true

tidy:  ## Tidy module dependencies.
	@$(foreach dir,$(MODULE_DIRS),( \
		cd $(dir) && \
		echo "tidy $(dir)" && \
		go mod tidy -compat=1.18) && \
		) true

test:  ## Run module tests.
	@$(foreach dir,$(MODULE_DIRS),( \
		cd $(dir) && \
		echo "test $(dir)" && \
		go test -race ./...) &&\
		) true

gen-backend: $(GQLGEN) ## Generate files for backend_service.
	@echo "generating GraphQL files for backend_service"
	@cd backend_service && $(GQLGEN)
	@echo "generating gRPC files for backend_service"
	@cd backend_service/proto/mongo && \
		protoc mongo.proto --go_out=. --go-grpc_out=.
	@cd backend_service/proto/postgres && \
		protoc postgres.proto --go_out=. --go-grpc_out=.

gen-mongo:  ## Generate files for mongo_service.
	@echo "generating gRPC files for mongo_service"
	@cd mongo_service/proto && \
		protoc mongo.proto --go_out=. --go-grpc_out=.

gen-postgres:  ## Generate files for postgres_service.
	@echo "generating gRPC files for postgres_service"
	@cd postgres_service/proto && \
		protoc postgres.proto --go_out=. --go-grpc_out=.

gen:  ## Generate files for all services.
	@make gen-backend
	@make gen-mongo
	@make gen-postgres

build:  ## Build development server.
	@echo "building docker compose network"
	@docker compose build

dev:  ## Start development server.
	@echo "starting docker compose network"
	@docker compose up -d

dev-build:  ## Build and start development server.
	@make build
	@make dev

clean:  ## Cleanup binary files.
	@echo "cleaning up 'bin'"
	@rm -rf bin/*

help: ## Display help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
