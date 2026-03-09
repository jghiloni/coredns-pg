GO ?= $(shell command -v go)

.PHONY: reset-generated generate test

generate: reset-generated
	TESTCONTAINERS_DOCKER_SOCKET_OVERRIDE=/var/run/docker.sock \
	TESTCONTAINERS_RYUK_DISABLED=true \
		$(GO) run ./codegen/ --use-temporary-db

reset-generated:
	rm -fr ./generated
test:
	$(GO) test ./...
