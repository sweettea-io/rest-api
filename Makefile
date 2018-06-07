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

run-server-go: ## Run the server application main file
	./scripts/run server

run-migrate-go: ## Run the migrate application main file
	./scripts/run migrate

run-worker-go: ## Run the work application main file
	./scripts/run worker

run-server-raw: ## Run the server application Go binary
	./scripts/run server -raw

run-migrate-raw: ## Run the migrate application Go binary
	./scripts/run migrate -raw

run-worker-raw: ## Run the worker application Go binary
	./scripts/run worker -raw

run-server: ## Run the latest Docker image for the server application
	./scripts/run server -docker

run-migrate: ## Run the latest Docker image for the migrate application
	./scripts/run migrate -docker

run-worker: ## Run the latest Docker image for the worker application
	./scripts/run worker -docker

run-dserver: ## Run the latest Docker image for the server application as a daemon
	./scripts/run server -docker -daemon

run-dmigrate: ## Run the latest Docker image for the migrate application as a daemon
	./scripts/run migrate -docker -daemon

run-dworker: ## Run the latest Docker image for the worker application as a daemon
	./scripts/run work -docker -daemon

deploy-server-local: ## Deploy the latest Docker image for the server application to the local K8S core cluster
	./scripts/deploy server local

deploy-server-dev: ## Deploy the latest Docker image for the server application to the dev K8S core cluster
	./scripts/deploy server dev

deploy-server-staging: ## Deploy the latest Docker image for the server application to the staging K8S core cluster
	./scripts/deploy server staging

deploy-server-prod: ## Deploy the latest Docker image for the server application to the prod K8S core cluster
	./scripts/deploy server prod

deploy-migrate-local: ## Deploy the latest Docker image for the migrate application to the local K8S core cluster
	./scripts/deploy migrate local

deploy-migrate-dev: ## Deploy the latest Docker image for the migrate application to the dev K8S core cluster
	./scripts/deploy migrate dev

deploy-migrate-staging: ## Deploy the latest Docker image for the migrate application to the staging K8S core cluster
	./scripts/deploy migrate staging

deploy-migrate-prod: ## Deploy the latest Docker image for the migrate application to the prod K8S core cluster
	./scripts/deploy migrate prod

deploy-worker-local: ## Deploy the latest Docker image for the worker application to the local K8S core cluster
	./scripts/deploy worker local

deploy-worker-dev: ## Deploy the latest Docker image for the worker application to the dev K8S core cluster
	./scripts/deploy worker dev

deploy-worker-staging: ## Deploy the latest Docker image for the worker application to the staging K8S core cluster
	./scripts/deploy worker staging

deploy-worker-prod: ## Deploy the latest Docker image for the worker application to the prod K8S core cluster
	./scripts/deploy worker prod

clean: ## Remove all Go-built binaries
	rm ./bin/*

ensure: ## Update dependencies
	dep ensure -vendor-only

test: ## Run all tests
	go test -p 1 ./...