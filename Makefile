.PHONY: build
PKG := "github.com/psyb0t/telegram-logger"
PKG_LIST := $(shell go list $(PKG)/...)

all: build

dep: ## Get the dependencies + remove unused ones
	@go mod tidy
	@go mod download

lint: ## Lint Golang files
	@golint -set_exit_status $(PKG_LIST)

build: dep ## Build the executable binary
	@go build -o build/app cmd/*.go

build-docker: ## Build the docker image via docker compose
	@docker compose build

run: dep ## Run without building
	@go run cmd/*.go

run-docker: ## Run in a docker container via docker compose
	@docker compose up

clean: ## Remove the build directory
	@rm -rf build

test: ## Run tests
	@go test -v $(PKG_LIST)

vet: ## Run go vet
	@go vet $(PKG_LIST)

test-coverage: ## Run tests with coverage
	@go test -short -coverprofile coverage.txt -covermode=atomic ${PKG_LIST}

test-coverage-tool: test-coverage ## Run test coverage followed by the cover tool
	@go tool cover -func=coverage.txt
	@go tool cover -html=coverage.txt

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
