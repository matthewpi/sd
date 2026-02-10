export CGO_ENABLED := '0'
pkg := 'github.com/matthewpi/sd'
goflags := '-v -trimpath'

build:
    @echo 'Building...'
    go build {{ goflags }} '{{ pkg }}/...'
    @echo 'Finished building!'

test:
    @echo 'Running tests...'
    go test {{ goflags }} '{{ pkg }}/...'

lint:
	@echo 'Linting project...'
	golangci-lint run --config .golangci.yaml

fmt:
	#!/bin/sh
	if [ -x "$(command -v nix)" ]; then
		nix fmt
		exit 0
	fi

	if [ -x "$(command -v gofumpt)" ]; then gofumpt -l -w .; fi
