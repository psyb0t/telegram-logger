.PHONY: build
MAKEFLAGS += --no-print-directory
PKG_BASENAME := "base-go-http-service"
PKG := "github.com/psyb0t/$(PKG)"
PKG_LIST := $(shell go list $(PKG)/...)

dep: ## Get the dependencies + remove unused ones
	@go mod tidy
	@go mod download

lint: ## Lint Golang files
	@golangci-lint run --timeout=30m0s

build: ## Build the executable binaries for all OSes and architectures
	@make build-linux-amd64
	@make build-windows-amd64
	@make build-darwin-amd64

build-linux-amd64: ## Build the executable binary for linux amd64
	@GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o build/$(PKG_BASENAME)-linux-amd64 cmd/*.go

build-windows-amd64: ## Build the executable binary for windows amd64
	@GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -installsuffix cgo -o build/$(PKG_BASENAME)-windows-amd64.exe cmd/*.go

build-darwin-amd64: ## Build the executable binary for darwin amd64
	@GO111MODULE=on CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -a -installsuffix cgo -o build/$(PKG_BASENAME)-darwin-amd64 cmd/*.go

build-docker: ## Build the docker image via docker compose
	@docker compose build

run: dep ## Run without building
	@go run -race cmd/*.go

run-docker: ## Run in a docker container via docker compose
	@docker compose up

clean: ## Remove the build directory
	@rm -rf build

vet: ## Run go vet
	@go vet $(PKG_LIST)

test: ## Run tests
	@go test -race -v $(PKG_LIST)

test-coverage: ## Run tests with coverage
	@go test -race -coverprofile coverage.txt -covermode=atomic ${PKG_LIST}

test-coverage-tool: test-coverage ## Run test coverage followed by the cover tool
	@go tool cover -func=coverage.txt
	@go tool cover -html=coverage.txt

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
