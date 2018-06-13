default: unlock

unlock: ## Unlock all scripts
	chmod -R u+x ./scripts/*

install: ## Setup & install project dependencies
	./scripts/install

export target
export env=local
export as=image
build: ## Build an application for a specific environment tier as either a Docker image or a Go binary.
	./scripts/build $(target) $(env) $(as)

export target
export as=image
export daemon
run: ## Run an application as either a Docker image, a Go binary, or a Go file.
	./scripts/run $(target) $(as) $(daemon)

export target
export env=local
deploy: ## Deploy the latest Docker image of an application to a specific environment.
	./scripts/deploy $(target) $(env)

export provider=aws
export type
export env
export zones
export master_size
export node_size
export node_count
export state
export image
export k8s_version
create-cluster: ## Create the core, train, or build cluster for a specific environment.
	./scripts/create-cluster $(provider) $(type) $(env) $(zones) $(master_size) $(node_size) $(node_count) $(state) $(image) $(k8s_version)

clean: ## Remove all Go-built binaries
	rm ./bin/*

ensure: ## Update dependencies
	dep ensure -vendor-only

test: ## Run all tests
	go test -p 1 ./...