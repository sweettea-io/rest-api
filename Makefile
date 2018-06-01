default: run

test: ## Run all tests
	go test -p 1 ./...

run: ## Start the HTTP server
	go run server.go

migrate: ## Run database migrations
	go run migrate.go

work: ## Start job worker
	go run worker.go

clean: ## Remove all built files
	rm server worker migrate

build: ## Build main executables
	go build ./...

build-server: rm server ## Build server.go
	go build -a -o server server.go

build-worker: rm worker ## Build worker.go
	go build -a -o worker worker.go

build-migrate: rm migrate ## Build migrate.go
	go build -a -o migrate migrate.go

update: # Update dependencies
	dep ensure