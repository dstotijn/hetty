# Ratatoskr build entry points. See docs/05-go-conventions.md §5.
# Recipes use POSIX sh. Supported platforms: Linux and macOS.

SHELL := /bin/sh

VERSION  ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS  := -s -w -X main.version=$(VERSION)
FUZZTIME ?= 10s
DATA_DIR ?= ./.dev-data

BIN := ratatoskr

# Web targets run only once the Svelte app exists (T-003); until then they are no-ops.
HAVE_WEB := $(wildcard web/package.json)

.PHONY: all
all: check build

## gen: regenerate code (gqlgen via go generate, web GraphQL codegen)
.PHONY: gen
gen: web-gen
	go generate ./...

## check: everything CI gates on
.PHONY: check
check: fix-check vet lint test web-check web-test

.PHONY: fix-check
fix-check:
	go fix -diff ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run ./...

# The race detector needs cgo on Linux, so tests opt back in.
.PHONY: test
test:
	CGO_ENABLED=1 go test -race ./...

## test-int: integration tests (build tag `integration`)
.PHONY: test-int
test-int:
	CGO_ENABLED=1 go test -race -tags integration ./...

## fuzz: run every Fuzz* target for FUZZTIME each
.PHONY: fuzz
fuzz:
	@for pkg in $$(go list ./...); do \
		for fn in $$(go test -list '^Fuzz' "$$pkg" | grep '^Fuzz'); do \
			echo "fuzz $$pkg $$fn"; \
			go test -run '^$$' -fuzz "^$$fn$$" -fuzztime $(FUZZTIME) "$$pkg" || exit 1; \
		done; \
	done

## vulncheck: scan dependencies for known vulnerabilities
.PHONY: vulncheck
vulncheck:
	go tool govulncheck ./...

## build: web bundle + single static binary
.PHONY: build
build: web-build
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN) ./cmd/ratatoskr

## run: build and start with a local data dir
.PHONY: run
run: build
	./$(BIN) --data-dir $(DATA_DIR)

.PHONY: web-gen web-check web-test web-build
web-gen:
ifneq ($(HAVE_WEB),)
	cd web && bun run codegen
else
	@echo "web-gen: skipped, no web/ app yet (T-003)"
endif

web-check:
ifneq ($(HAVE_WEB),)
	cd web && bun run check
else
	@echo "web-check: skipped, no web/ app yet (T-003)"
endif

web-test:
ifneq ($(HAVE_WEB),)
	cd web && bun run test
else
	@echo "web-test: skipped, no web/ app yet (T-003)"
endif

web-build:
ifneq ($(HAVE_WEB),)
	cd web && bun run build
else
	@echo "web-build: skipped, no web/ app yet (T-003)"
endif

.PHONY: clean
clean:
	rm -f $(BIN)
