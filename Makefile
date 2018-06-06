default: setup

setup: ## Unlock all scripts
	chmod u+x scripts/*

install: ## Install project dependencies
	./scripts/install

test: ## Run all tests
	go test -p 1 ./...

run: ## Start the HTTP server
	go run server.go

migrate: ## Run database migrations
	go run migrate.go

work: ## Start job worker
	go run worker.go

clean: ## Remove all built files
	rm bin/*

build-server: ## Build server.go
	go build -a -o ./bin/server ./cmd/server

build-worker: ## Build worker.go
	go build -a -o ./bin/worker ./cmd/worker

build-migrate: ## Build migrate.go
	go build -a -o ./bin/migrate ./cmd/migrate

ensure: # Update dependencies
	dep ensure -vendor-only