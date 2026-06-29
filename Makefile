export CGO_ENABLED = 0
export NEXT_TELEMETRY_DISABLED = 1

.PHONY: build
build: build-admin
	go build ./cmd/hetty

.PHONY: build-desktop
build-desktop: build-admin
	cd cmd/hetty-desktop && CGO_ENABLED=1 GOTOOLCHAIN=go1.25.11 go build -o ../../hetty-desktop .

.PHONY: build-admin
build-admin:
	cd admin && \
	npx --yes yarn install --frozen-lockfile && \
	npx yarn run export && \
	mv dist ../cmd/hetty/admin && \
	cp -r ../cmd/hetty/admin ../cmd/hetty-desktop/admin

.PHONY: clean
clean:
	rm -f hetty hetty-desktop
	rm -rf ./cmd/hetty/admin
	rm -rf ./cmd/hetty-desktop/admin
	rm -rf ./admin/dist
	rm -rf ./admin/.next