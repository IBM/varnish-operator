# Monitoring

## Operator monitoring

The operator is built using the [Kubebuilder SDK](https://github.com/kubernetes-sigs/kubebuilder) which has built-in support for the Prometheus metrics exporter.

By default, the operator creates a Service that exposes the port the exporter is listening on. The default port is `9131` and can be changed by overriding the `container.metricsPort` field in [Helm chart](operator-configuration.md).

## Varnish monitoring

Each Varnish pod has a [Varnish Prometheus metrics exporter](https://github.com/jonnenauha/prometheus_varnish_exporter) built-in. They exporter port is exposed by the `VarnishService` on port `9131` by default. It can be change by setting the `spec.service.varnishExporterPort.port` field in the [`VarnishService` spec](varnish-service-configuration.md).

The service port can be used to setup metrics scraping using [Prometheus Operator](https://github.com/coreos/prometheus-operator) `ServiceMonitor`.  

The pods itself expose metrics on port `9131`. 

