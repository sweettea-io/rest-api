default: unlock

unlock: ## Unlock all scripts.
	chmod -R +x ./scripts/*

install: ## Install project dependencies.
	./scripts/install

export password
user-creation-hash: ## Create the user-creation hash to be used on the environments of your choice.
	go run tasks/generate_user_creation_hash.go $(password)

local-clusters: ## Create the core, train, and build clusters for the local environment.
	./scripts/create_cluster '' core local
	./scripts/create_cluster '' train local
	./scripts/create_cluster '' build local

export cloud=aws
export role
export env=local
cluster: ## Create the core, train, or build cluster for a specific environment.
	./scripts/create_cluster $(cloud) $(role) $(env)

export cloud=aws
export env=local
model-storage: ## Create a new model storage instance to use with a specific environment.
	./scripts/create_model_storage $(cloud) $(env)

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

export cloud=aws
export name
export env=local
api: ## Create a new API cluster for a specific environment.
	./scripts/create_api $(cloud) $(name) $(env)

export env=local
export email
export password
export admin=false
user: ## Create Sweet Tea User in the database of the specified environment.
	./scripts/create_user $(env) $(email) $(password) $(admin)

clean: ## Remove all built Go binaries.
	rm ./bin/*

ensure: ## Update dependencies.
	./scripts/install_pkgs

test: ## Run all tests.
	./scripts/run_tests

rebuild-local-clusters: ## Delete and re-create the core, train, and build clusters for the local environment.
	rm -rf ~/.minikube
	./scripts/rebuild_local_cluster core
	./scripts/rebuild_local_cluster train
	./scripts/rebuild_local_cluster build