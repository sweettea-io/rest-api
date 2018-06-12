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

clean: ## Remove all Go-built binaries
	rm ./bin/*

ensure: ## Update dependencies
	dep ensure -vendor-only

test: ## Run all tests
	go test -p 1 ./...