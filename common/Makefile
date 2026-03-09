GO ?= $(shell command -v go)
GOLANGCI_LINT ?= $(shell command -v golangci-lint)

CITARGETS = run fmt

.PHONY: reset-generated generate ci-run ci-fmt

generate: reset-generated
	TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE=/var/run/docker.sock \
	TESTCONTAINERS_RYUK_DISABLED=true \
		$(GO) run ./codegen/ --use-temporary-db

reset-generated:
	rm -fr ./generated

test:
	$(GO) test ./...

$(CITARGETS):
	$(GOLANGCI_LINT) $@
