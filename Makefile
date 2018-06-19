default: unlock

unlock: ## Unlock all scripts
	chmod -R u+x ./scripts/*

install: ## Setup & install project dependencies
	./scripts/install

export provider=aws
export type
export env
cluster: ## Create the core, train, or build cluster for a specific environment.
	./scripts/create_cluster $(provider) $(type) $(env)

export provider=aws
export env
clusters: ## Make all clusters for a specific environment.
	./scripts/create_cluster $(provider) core $(env)
	./scripts/create_cluster $(provider) build $(env)
	./scripts/create_cluster $(provider) train $(env)

export target
export env=local
export format=image
build: ## Build an application for a specific environment tier as either a Docker image or a Go binary.
	./scripts/build $(target) $(env) $(format)

export target
export format=image
export daemon
run: ## Run an application as either a Docker image, a Go binary, or a Go file.
	./scripts/run $(target) $(format) $(daemon)

export target
export env=local
deploy: ## Deploy the latest Docker image of an application to a specific environment.
	./scripts/deploy $(target) $(env)

clean: ## Remove all built Go binaries
	rm ./bin/*

ensure: ## Update dependencies
	dep ensure -vendor-only

test: ## Run all tests
	go test -p 1 ./...