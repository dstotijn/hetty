PACKAGE_NAME := github.com/dstotijn/hetty
GOLANG_CROSS_VERSION ?= v1.16.3

.PHONY: build-admin
build-admin:
	NEXT_TELEMETRY_DISABLED=1 cd admin && yarn install && yarn run export

.PHONY: build
build: build-admin
	CGO_ENABLED=1 mv admin/dist cmd/hetty/admin && go build ./cmd/hetty

.PHONY: release-dry-run
release-dry-run: build-admin
	@docker run \
		--rm \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		troian/golang-cross:${GOLANG_CROSS_VERSION} \
		--rm-dist --skip-validate --skip-publish

.PHONY: release
release: build-admin
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91mFile \`.release-env\` is missing.\033[0m";\
		exit 1;\
	fi
	@docker run \
		--rm \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		--env-file .release-env \
		troian/golang-cross:${GOLANG_CROSS_VERSION} \
		release --rm-dist