default: unlock

unlock: ## Unlock all scripts
	chmod -R +x ./scripts/*

install: ## Install project dependencies
	./scripts/install

export password
user-creation-hash: ## Create the user-creation hash to be used on the environments of your choice.
	go run tasks/generate_user_creation_hash.go $(password)

local-clusters: ## Create the core, train, and build clusters for the local environment.
	./scripts/create_cluster '' core local
	./scripts/create_cluster '' train local
	./scripts/create_cluster '' build local

export provider=aws
export role
export env=local
cluster: ## Create the core, train, or build cluster for a specific environment.
	./scripts/create_cluster $(provider) $(role) $(env)

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

export provider=aws
export name
export company
export env=local
api: ## Create a new API cluster for a specific environment.
	./scripts/create_api $(provider) $(name) $(company) $(env)

export env=local
export email
export password
export admin=false
user: ## Create Sweet Tea User in the database of the specified environment.
	./scripts/create_user $(env) $(email) $(password) $(admin)

export env=local
export name
export with-api=false
company: ## Create Sweet Tea Company in the database of the specified environment.
	./scripts/create_company $(env) $(name)

clean: ## Remove all built Go binaries
	rm ./bin/*

ensure: ## Update dependencies
	dep ensure -vendor-only

test: ## Run all tests
	go test -p 1 ./...