default: unlock

unlock: ## Unlock all scripts
	chmod -R u+x scripts/*

install: ## Setup & install project dependencies
	./scripts/install

build-server: ## Build the server application as a Docker image
	./scripts/build server

build-migrate: ## Build the migrate application as a Docker image
	./scripts/build migrate

build-worker: ## Build the worker application as a Docker image
	./scripts/build worker

build: ## Build all applications as Docker images
	./scripts/build all

build-server-raw: ## Build the server application as a Go binary
	./scripts/build server -raw

build-migrate-raw: ## Build the migrate application as a Go binary
	./scripts/build migrate -raw

build-worker-raw: ## Build the worker application as a Go binary
	./scripts/build worker -raw

build-raw: ## Build all applications as Go binaries
	./scripts/build all -raw

run: ## Start the HTTP server
	go run server.go

migrate: ## Run database migrations
	go run migrate.go

work: ## Start job worker
	go run worker.go

clean: ## Remove all built files
	rm bin/*

ensure: ## Update dependencies
	dep ensure -vendor-only

test: ## Run all tests
	go test -p 1 ./...

