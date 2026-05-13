export CGO_ENABLED = 0
export NEXT_TELEMETRY_DISABLED = 1

.PHONY: build
build: build-admin
	go build ./cmd/hetty

.PHONY: build-admin
build-admin:
	cd admin && \
	yarn install --frozen-lockfile && \
	yarn run export && \
	rm -rf ../cmd/hetty/admin && \
	cp -R dist ../cmd/hetty/admin

.PHONY: clean
clean:
	rm -f hetty
	rm -rf ./cmd/hetty/admin
	rm -rf ./admin/dist
	rm -rf ./admin/.next