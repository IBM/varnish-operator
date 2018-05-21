FROM alpine:3.6

ADD tmp/_output/bin/icm-varnish-k8s-operator /usr/local/bin/icm-varnish-k8s-operator

RUN adduser -D icm-varnish-k8s-operator
USER icm-varnish-k8s-operator
