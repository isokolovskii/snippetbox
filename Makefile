.PHONY: run
run:
	go tool air -c .air.toml

.PHONY: build
build:
	mkdir -p tmp
	go build -o tmp/web ./cmd/web

.PHONY: build-with-coverage
build-with-coverage:
	go build -cover ./...

.PHONY: test
test:
	go test -cpu 24 -race -count=1 -timeout=30s ./...

.PHONY: bench
bench:
	go test -cpu 24 -race -run=Bench -bench=. ./...

GOLANGCI_LINT_BIN := $(shell go env GOPATH)/bin/golangci-lint

.PHONY: lint
lint: $(GOLANGCI_LINT_BIN)
	$(GOLANGCI_LINT_BIN) run ./...

.PHONY: format
format: $(GOLANGCI_LINT_BIN)
	$(GOLANGCI_LINT_BIN) fmt ./...

$(GOLANGCI_LINT_BIN):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.6.2

.PHONY: lefthook
lefthook:
	go tool lefthook install

.PHONY: create-certs
create-certs:
	mkdir -p tls
	cd tls; \
	go run $(shell go env GOROOT)/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost

.PHONY: init
init: lefthook create-certs

.PHONY: create-migration
create-migration:
	@read -p "Enter migration name:" migration; \
		$(shell go env GOPATH)/bin/migrate create -ext sql -dir migrations -seq $$migration

.PHONY: compose-up
compose-up:
	docker compose up --watch

.PHONY: compose-build
compose-build:
	docker compose up --watch --build

.PHONY: compose-down
compose-down:
	docker compose down

.PHONY: clear-local-db
clear-local-db:
	rm -rf db/data

.PHONY: recreate-local-db
recreate-local-db: compose-down clear-local-db compose-up
