default: unlock

unlock: ## Unlock all scripts
	chmod -R +x ./scripts/*

install: ## Setup & install project dependencies
	./scripts/install

export provider=aws
export role
export env
cluster: ## Create the core, train, or build cluster for a specific environment.
	./scripts/create_cluster $(provider) $(role) $(env)

export provider=aws
export name
export env
api: ## Create a new API cluster for a specific environment.
	./scripts/create_api $(provider) $(name) $(env)

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

export env
export email
export pw
export admin=false
user: ## Create Sweet Tea User in the database of the specified environment.
	./scripts/create_user $(env) $(email) $(pw) $(admin)

clean: ## Remove all built Go binaries
	rm ./bin/*

ensure: ## Update dependencies
	dep ensure -vendor-only

test: ## Run all tests
	go test -p 1 ./...