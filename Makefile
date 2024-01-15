# Tool versions
99_VERSION := v0.17.40
GOLANGCI_LINT_VERSION := v1.55.2

.PHONY: setup-mac
setup-mac: ## Setup tool ecosystem for MacOS
setup-mac: setup-mac-install-asdf setup-common-install

.PHONY: setup-mac-install-asdf
setup-mac-install-asdf:
	brew install asdf

	# Idempotently add language plugins.
	asdf plugin-add golang || true
	asdf plugin-add python || true
	asdf plugin-add nodejs || true
	asdf plugin-add yarn || true

	asdf install

	# Ensure that asdf is in the path... 
	@(which go | grep asdf -q) || \
	  (echo "ERROR: asdf not in PATH. You're using the wrong tools ;). Follow instructions for adding to shell after a brew install." && \
	    open "https://asdf-vm.com/guide/getting-started.html#_3-install-asdf" && \
		false)

.PHONY: setup-common-install
setup-common-install: setup-common-gen

.PHONY: setup-common-gen
setup-common-gen: 
	go install github.com/99designs/gqlgen@$(99_VERSION)

.PHONY: test-unit-go
test-unit-go: gen-gql
	go test ./...

.PHONY: gen-gql
gen-gql:
	go run github.com/99designs/gqlgen generate --config $(CURDIR)/gqlgen.yml

.PHONY: run-go-backend
run-go-backend: gen-gql
	go run go/main.go

.PHONY: run-client
run-client:
	yarn
	yarn dev


