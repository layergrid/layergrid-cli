VERSION ?= 0.1.0-dev
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)

.PHONY: build test lint fmt install bench principles docs scan-corpus bench-realworld validate-release

build:
	go build -ldflags "-s -w -X github.com/layergrid/layergrid-cli/internal/version.Version=$(VERSION) -X github.com/layergrid/layergrid-cli/internal/version.Commit=$(COMMIT) -X github.com/layergrid/layergrid-cli/internal/version.Date=$(DATE)" -o bin/layergrid ./cmd/layergrid

test:
	go test -race -cover ./...

bench:
	go test -bench=. -benchmem ./internal/scan

scan-corpus: build
	./hack/scan-corpus.sh

bench-realworld: build
	./hack/bench-realworld.sh

lint:
	golangci-lint run

principles:
	./scripts/principles-check.sh

docs:
	./hack/generate-rules-doc.sh

validate-release:
	./scripts/validate-release.sh

fmt:
	gofmt -s -w .

install: build
	install -m 0755 bin/layergrid $(GOBIN)/layergrid
