default: run

install: ## Setup project
	mkdir bin
	mkdir envs
	update

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

update: # Update dependencies
	dep ensure