# Build the etcd-certs-to-secret binary
# Build context must be the repository root.
FROM --platform=$BUILDPLATFORM golang:1.27.0-alpine3.24@sha256:4c9fe60190a2a3350ddc51de80d0224b8a6698d12bdfc999fee45ea9d6c46dbc AS builder
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
FROM alpine:3.24.1@sha256:28bd5fe8b56d1bd048e5babf5b10710ebe0bae67db86916198a6eec434943f8b

ENV USER_UID=2001 \
    USER_NAME=etcd-certs \
    GROUP_NAME=etcd-certs

WORKDIR /
COPY --from=builder --chown=${USER_UID} --chmod=0555 /workspace/etcd-certs .

RUN addgroup "${GROUP_NAME}" && adduser -D -G "${GROUP_NAME}" -u "${USER_UID}" "${USER_NAME}"
USER ${USER_UID}

ENTRYPOINT ["/etcd-certs"]
