setup:
	go mod download
	go generate ./...
.PHONY: setup

embed:
	go install github.com/GeertJohan/go.rice/rice
	cd cmd/hetty && rice embed-go
.PHONY: embed

build: embed
	go build ./cmd/hetty
.PHONY: build

clean:
	rm -rf cmd/hetty/rice-box.go
.PHONY: clean

release:
	goreleaser -p 1