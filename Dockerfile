ARG GO_VERSION=1.16
ARG NODE_VERSION=14.11
ARG ALPINE_VERSION=3.12

ARG CGO_ENABLED=1

FROM node:${NODE_VERSION}-alpine AS node-builder
WORKDIR /app
COPY admin/package.json admin/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY admin/ .
ENV NEXT_TELEMETRY_DISABLED=1
RUN yarn run export

FROM golang:${GO_VERSION}-alpine AS go-builder
WORKDIR /app
RUN apk add --no-cache build-base
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY pkg ./pkg
COPY --from=node-builder /app/dist ./cmd/hetty/admin
RUN go build ./cmd/hetty

FROM alpine:${ALPINE_VERSION}
WORKDIR /app
COPY --from=go-builder /app/hetty .

ENTRYPOINT ["./hetty"]

EXPOSE 8080