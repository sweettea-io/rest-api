default: unlock

unlock: ## Unlock all scripts
	chmod -R u+x ./scripts/*

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

serve-go: ## Run the server application main file
	./scripts/run server

migrate-go: ## Run the migrate application main file
	./scripts/run migrate

work-go: ## Run the work application main file
	./scripts/run worker

serve-raw: ## Run the server application Go binary
	./scripts/run server -raw

migrate-raw: ## Run the migrate application Go binary
	./scripts/run migrate -raw

work-raw: ## Run the worker application Go binary
	./scripts/run worker -raw

serve: ## Run the latest Docker image for the server application
	./scripts/run server -docker

migrate: ## Run the latest Docker image for the migrate application
	./scripts/run migrate -docker

work: ## Run the latest Docker image for the worker application
	./scripts/run work -docker

dserve: ## Run the latest Docker image for the server application as a daemon
	./scripts/run server -docker -daemon

dmigrate: ## Run the latest Docker image for the migrate application as a daemon
	./scripts/run migrate -docker -daemon

dwork: ## Run the latest Docker image for the worker application as a daemon
	./scripts/run work -docker -daemon

clean: ## Remove all Go-built binaries
	rm ./bin/*

ensure: ## Update dependencies
	dep ensure -vendor-only

test: ## Run all tests
	go test -p 1 ./...
