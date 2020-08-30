ARG GO_VERSION=1.14.4-alpine
FROM golang:${GO_VERSION} as builder

ARG CGO_ENABLED=0

WORKDIR /go/src/myapp
COPY go.mod go.sum ./
RUN set -eux; \
    go mod download -x

# Copy the local package files to the container's workspace.
COPY . ./

# Run unittests first
# RUN set -ux; \
#     go test -cover ./...

ARG GOOS=linux
ARG GOARCH=amd64
ARG RAW_VERSION
ARG BUILD_DATE
ARG GIT_COMMIT
ARG GIT_TAG
ARG GIT_TREE_STATE
ARG PACKAGE=github.com/nikser/simple-k8s-app
ARG GO_LDFLAGS="-w -s -d"
ARG GO_XFLAGS="-X \"${PACKAGE}.version=${RAW_VERSION}\" -X \"${PACKAGE}.buildDate=${BUILD_DATE}\" -X \"${PACKAGE}.gitCommit=${GIT_COMMIT}\" -X \"${PACKAGE}.gitTreeState=${GIT_TREE_STATE}\" -X \"${PACKAGE}.gitTag=${GIT_TAG}\""

# Support for multiple binaries
RUN set -eux; \
    go build -ldflags "${GO_LDFLAGS} ${GO_XFLAGS}" -o /go/bin/myapp .

# This results in a single layer image
FROM scratch

COPY --from=builder /go/bin/myapp /

ENTRYPOINT ["/myapp"]
