ARG GO_VERSION=1.15
ARG CGO_ENABLED=1
ARG NODE_VERSION=14.11

FROM golang:${GO_VERSION} AS go-builder
WORKDIR /app
RUN apt-get update && \
    apt-get install -y build-essential
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY pkg ./pkg
ENV CGO_CFLAGS=-I/go/pkg/mod/github.com/mattn/go-sqlite3@v1.14.4
ENV CGO_LDFLAGS=-Wl,--unresolved-symbols=ignore-in-object-files
RUN go build -o hetty ./cmd/hetty

FROM node:${NODE_VERSION}-alpine AS node-builder
WORKDIR /app
COPY admin/package.json admin/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY admin/ .
ENV NEXT_TELEMETRY_DISABLED=1
RUN yarn run export

FROM debian:buster-slim
WORKDIR /app
COPY --from=go-builder /app/hetty .
COPY --from=node-builder /app/dist admin

ENTRYPOINT ["./hetty", "-adminPath=./admin"]

EXPOSE 8080