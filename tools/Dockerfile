ARG GO_VERSION

FROM golang:${GO_VERSION}-alpine

ARG GOLANGCI_VERSION=v1.41.1
ARG GOFUMPT_VERSION=v0.1.1
ARG GO_IMPORT_VERSION=v0.1.4

RUN apk update && apk add --no-cache \
    curl \
    docker \
    git \
    unzip

# binary will be $(go env GOPATH)/bin/golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_VERSION}

# Install Go packages
RUN go install mvdan.cc/gofumpt@${GOFUMPT_VERSION} 
RUN go install golang.org/x/tools/cmd/goimports@${GO_IMPORT_VERSION}

# Set default golangci config
COPY .golangci.yml /.golangci.yml

