# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.25.5-alpine3.22 AS builder
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

# Build
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} GO111MODULE=on go build -a -o manager main.go

# Use alpine base container
FROM alpine:3.23.2

ENV USER_UID=2001 \
    USER_NAME=monitoring-operator \
    GROUP_NAME=monitoring-operator

WORKDIR /

# Create user and group first
RUN addgroup ${GROUP_NAME} && adduser -D -G ${GROUP_NAME} -u ${USER_UID} ${USER_NAME}

# Copy manager binary from builder stage with correct ownership
COPY --from=builder --chown=${USER_UID}:${USER_UID} /workspace/manager /manager

# Ensure the binary is executable but not writable (security best practice)
# Set permissions to read and execute only (555 = r-xr-xr-x)
RUN chmod 555 /manager

USER ${USER_UID}

ENTRYPOINT ["/manager"]
