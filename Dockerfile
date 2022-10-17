ARG BUILDPLATFORM=linux/amd64
FROM --platform=$BUILDPLATFORM golang:1.19.2-bullseye AS builder

ENV DEBIAN_FRONTEND=noninteractive INSTALL_DIRECTORY=/usr/local/bin

RUN apt-get update && apt-get install -y --no-install-recommends \
        git \
        curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/github.com/ibm/varnish-operator

ENV GOPROXY=https://proxy.golang.org

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/
COPY api/ api/

ARG TARGETOS=linux
ARG TARGETARCH=amd64
ARG VERSION=undefined

# Build
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build \
    -ldflags "-X main.Version=$VERSION" \
    -a \
    -o varnish-operator \
    github.com/ibm/varnish-operator/cmd/varnish-operator


FROM --platform=$BUILDPLATFORM debian:bullseye-slim

LABEL maintainer="Alex Lytvynenko <oleksandr.lytvynenko@ibm.com>, Tomash Sidei <tomash.sidei@ibm.com>"

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /

RUN apt-get update && apt-get upgrade -y \
    && rm -rf /var/lib/apt/lists/*

RUN addgroup --gid 901 varnish-operator && adduser --uid 901 --gid 901 varnish-operator

COPY --from=builder /go/src/github.com/ibm/varnish-operator/varnish-operator .

USER varnish-operator

ENTRYPOINT ["/varnish-operator"]
