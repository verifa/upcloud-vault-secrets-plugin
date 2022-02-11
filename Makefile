.DEFAULT_GOAL := all

all: fmt build start

build:
	go build -o build/upcloud-vault-secrets-plugin cmd/main.go

start:
	vault server -dev -dev-root-token-id=root -dev-plugin-dir=./build -log-level=trace

enable:
	vault secrets enable -path=upcloud upcloud-vault-secrets-plugin

clean:
	rm -f ./build/upcloud-vault-secrets-plugin

fmt:
	go fmt $$(go list ./...)

.PHONY: build clean fmt start enable
