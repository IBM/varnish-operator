# POC for Using Varnish in Kubernetes

## K-Watcher

Go app that watches kubernetes endpoints api for changes in deployment, and then re-writes varnish vcl file with any new/removed backends.

## Docker

Docker image with Varnish cache 5.x software and k-watcher packaged in Alpine Linux docker image.

