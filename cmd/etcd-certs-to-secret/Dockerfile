# Build the etcd-certs-to-secret binary
# Build context must be the repository root.
FROM --platform=$BUILDPLATFORM golang:1.26.3-alpine3.23 AS builder
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Cache deps before copying source
RUN go mod download

# Copy the go source
COPY cmd/etcd-certs-to-secret/ cmd/etcd-certs-to-secret/

# Build
RUN CGO_ENABLED=0 GOOS="${TARGETOS}" GOARCH="${TARGETARCH}" GO111MODULE=on \
    go build -a -o etcd-certs ./cmd/etcd-certs-to-secret

# Use alpine base container
FROM alpine:3.23.4

ENV USER_UID=2001 \
    USER_NAME=etcd-certs \
    GROUP_NAME=etcd-certs

WORKDIR /
COPY --from=builder --chown=${USER_UID} --chmod=0555 /workspace/etcd-certs .

RUN addgroup "${GROUP_NAME}" && adduser -D -G "${GROUP_NAME}" -u "${USER_UID}" "${USER_NAME}"
USER ${USER_UID}

ENTRYPOINT ["/etcd-certs"]
