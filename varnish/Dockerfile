FROM billyteves/alpine-golang-glide:1.2.0 AS builder

WORKDIR /go/src/icm-varnish/varnish/k-watcher
COPY k-watcher/cmd cmd
COPY k-watcher/main.go k-watcher/glide.yaml ./
RUN set -ex &&\
    glide update &&\
    glide install -v
RUN go build -o k-watcher .

FROM alpine:3.7
LABEL maintainer="Alex Lytvynenko <oleksandr.lytvynenko@ibm.com>"

ENV PROMETHEUS_VARNISH_EXPORTER_VERSION 1.4

ENV VARNISH_MEMORY 2048M
ENV VARNISH_PORT 8080
# this should be moved to a build-arg
ENV VARNISH_SECRET c5m4HSNJwFz7A1Hygc0xoENSsoEFwTr

RUN apk update && \
    apk upgrade && \
    apk add varnish curl python py-pip && \
    easy_install-2.7 supervisor && \
    curl -L https://github.com/jonnenauha/prometheus_varnish_exporter/releases/download/${PROMETHEUS_VARNISH_EXPORTER_VERSION}/prometheus_varnish_exporter-${PROMETHEUS_VARNISH_EXPORTER_VERSION}.linux-amd64.tar.gz \
    | tar --strip-components=1 -xz -C /usr/bin && \
    chmod +x /usr/bin/prometheus_varnish_exporter && \
    mkdir /lib64 && \
    ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \
    apk del curl && \
    rm -vf /var/cache/apk/*

COPY --from=builder /go/src/icm-varnish/varnish/k-watcher/k-watcher /k-watcher
COPY supervisord.conf /etc/supervisor/supervisord.conf
COPY start_varnishd.sh start_k_watcher.sh /
COPY vcl_reload /usr/local/bin/
COPY varnish/config/default.vcl varnish/config/backends.vcl k-watcher/backends.vcl.tmpl /etc/varnish/

RUN chown -R varnish /etc/varnish

USER varnish

CMD ["supervisord", "-c", "/etc/supervisor/supervisord.conf"]
