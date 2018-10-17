FROM golang:1.11.1-alpine3.8 AS builder

RUN apk update && apk add curl git

ENV DEP_RELEASE_TAG=v0.5.0 INSTALL_DIRECTORY=/usr/local/bin
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/icm-varnish-k8s-operator

# Copy in the go src
COPY Gopkg.toml Gopkg.lock ./
COPY cmd/       cmd/
COPY pkg/       pkg/

# Populate the vendor folder
RUN dep ensure -v

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager icm-varnish-k8s-operator/cmd/manager

FROM alpine:3.8
LABEL maintainer="thurston sandberg <thurston.sandberg@us.ibm.com>"

WORKDIR /

RUN apk update &&\
    apk upgrade

RUN addgroup -g 901 controller && adduser -D -u 901 -G controller controller
# RUN chown -R controller

COPY config/vcl/default.vcl config/vcl/backends.vcl.tmpl config/vcl/
COPY --from=builder /go/src/icm-varnish-k8s-operator/manager .

USER controller
ENTRYPOINT ["/manager"]
