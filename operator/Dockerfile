FROM billyteves/alpine-golang-glide:1.2.0 AS builder

WORKDIR /go/src/icm-varnish/operator/controller
COPY controller/cmd cmd
COPY controller/main.go controller/glide.yaml ./
RUN set -ex &&\
  glide update &&\
  glide install -v
RUN go build -o varnish-service-controller .

FROM alpine:3.7
LABEL maintainer="thurston sandberg <thurston.sandberg@us.ibm.com>"

RUN apk update &&\
  apk upgrade

COPY --from=builder /go/src/icm-varnish/operator/controller/varnish-service-controller /varnish-service-controller
COPY start_controller.sh /

USER controller
# RUN chown -R controller 
CMD ["/start_controller.sh"]