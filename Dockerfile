FROM golang:1.14-buster AS builder

ARG GOPROXY=direct

ENV DEBIAN_FRONTEND=noninteractive INSTALL_DIRECTORY=/usr/local/bin

RUN apt-get update && apt-get install -y --no-install-recommends \
        git \
        curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/icm-varnish-k8s-operator

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY pkg/ pkg/
COPY api/ api/
COPY version.txt ./

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags "-X main.Version=$(cat ./version.txt)" \
    -a \
    -o varnish-operator \
    icm-varnish-k8s-operator/cmd/varnish-operator


FROM debian:buster-slim

LABEL maintainer="Alex Lytvynenko <oleksandr.lytvynenko@ibm.com>, Tomash Sidei <tomash.sidei@ibm.com>"

ENV DEBIAN_FRONTEND=noninteractive

WORKDIR /

RUN apt-get update && apt-get upgrade -y \
    && rm -rf /var/lib/apt/lists/*

RUN addgroup --gid 901 varnish-operator && adduser --uid 901 --gid 901 varnish-operator

COPY --from=builder /go/src/icm-varnish-k8s-operator/varnish-operator .

USER varnish-operator

ENTRYPOINT ["/varnish-operator"]
