export CGO_ENABLED = 0
export NEXT_TELEMETRY_DISABLED = 1

.PHONY: clean
clean:
	rm -f hetty
	rm -rf ./cmd/hetty/admin
	rm -rf ./admin/node_modules
	rm -rf ./admin/dist
	rm -rf ./admin/.next

.PHONY: build-admin
build-admin:
	cd admin && \
	yarn install --frozen-lockfile && \
	yarn run export && \
    mv dist ../cmd/hetty/admin

.PHONY: build
build: build-admin
	go build ./cmd/hetty