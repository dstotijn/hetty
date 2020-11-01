ARG GO_VERSION=1.15
ARG CGO_ENABLED=1
ARG NODE_VERSION=14.11

FROM golang:${GO_VERSION}-alpine AS go-builder
WORKDIR /app
RUN apk add --no-cache build-base
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY pkg ./pkg
RUN rm -f cmd/hetty/rice-box.go
RUN go build ./cmd/hetty

FROM node:${NODE_VERSION}-alpine AS node-builder
WORKDIR /app
COPY admin/package.json admin/yarn.lock ./
RUN yarn install --frozen-lockfile
COPY admin/ .
ENV NEXT_TELEMETRY_DISABLED=1
RUN yarn run export

FROM alpine:3.12
WORKDIR /app
COPY --from=go-builder /app/hetty .
COPY --from=node-builder /app/dist admin

ENTRYPOINT ["./hetty", "-adminPath=./admin"]

EXPOSE 8080