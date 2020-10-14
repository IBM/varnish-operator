# Monitoring

## Operator Monitoring

The operator is built using the [Kubebuilder SDK](https://github.com/kubernetes-sigs/kubebuilder) which has built-in support for the Prometheus metrics exporter.

A service, created by the operator's Helm chart, exposes the metrics on port `8329` (named `prometheus-metrics`) and can be used to scrape operator metrics.

Additionally, the operator can install a ServiceMonitor configured to scrape operator metrics and a Grafana dashboard with prebuilt dashboard for the operator. The configuration options for them can be specified under `.monitoring` [values override](operator-configuration.md) field of the Helm chart.

### Monitoring Stack Example

The repo includes an [example helm chart](https://github.com/IBM/varnish-operator/tree/main/config/samples/helm-charts/varnish-operator-monitoring) for a Prometheus and Grafana installation that is configured to scrape metrics from the operator and display included dashboards. It depends on the [Prometheus Operator](https://github.com/coreos/prometheus-operator) so it has to be installed first.

After you have the [operator installed](installation.md), clone the repo and install the helm chart.

```bash
$ git clone https://github.com/IBM/varnish-operator.git
$ cd varnish-operator/config/samples/helm-charts/varnish-operator-monitoring
$ helm dep build
$ helm install --name varnish-operator-monitoring .
```

No additional configuration needed. The monitoring stack relies on the labels set for the Service that exposes the operator pods.

Port forward your grafana installation:

```bash
$ kubectl port-forward pod/varnish-operator-monitoring-grafana-6f7ff7f4f9-2pjpj 3000
Forwarding from 127.0.0.1:3000 -> 3000
Forwarding from [::1]:3000 -> 3000
```

You can see your dashboard at `localhost:3000`. The login is `admin`, password is `pass`. You will find a dashboard called `Varnish Operator`.

The chart is not a complete solution and intended to be modified to the end user needs.

## Varnish Monitoring

Each Varnish pod has a [Varnish Prometheus metrics exporter](https://github.com/jonnenauha/prometheus_varnish_exporter) built-in. They exporter port is exposed by the `VarnishCluster` on port `9131` by default. It can be changed by setting the `spec.service.prometheusExporterPort` field in the [`VarnishCluster` spec](varnish-cluster-configuration.md).

The service port can be used to setup metrics scraping using [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator) `ServiceMonitor`.  

The pods itself also expose metrics on port `9131`.

There are also metrics exposed by the Varnish controller on a different port. You'll have to setup your ServiceMonitor to scrape metrics from the port 8235 (or refer by its name `ctrl-metrics` 

One of the metrics can be used to setup alerts when the provided VCL failed to compile. It's called `varnish_vcl_compilation_error` and has the value `0` if the last compilation attempt was successful or `1` in case of failure.

### VarnishCluster with Monitoring Stack Example

The repo has a Helm chart example that installs a simple VarnishCluster with Prometheus and Grafana configured to monitor it. The chart depends on the Prometheus operator so it has to be installed first. 

It can be installed by cloning the repo and installing the chart with necessary backend configs:

```bash
$ git clone https://github.com/IBM/varnish-operator.git
$ cd varnish-operator/config/samples/helm-charts/varnishcluster-with-monitoring
$ helm dep build
$ helm install --name varnish-test . --set varnish.backendsSelector.app=nginx --set varnish.backendsPort=80
```

Port forward the Grafana pod to see the dashboard:

```bash
$ kubectl port-forward pod/varnish-test-grafana-9f584598d-89smp 3000
Forwarding from 127.0.0.1:3000 -> 3000
Forwarding from [::1]:3000 -> 3000
```

You can see your dashboard at `localhost:3000`. The login is `admin`, password is `pass`. You will find a dashboard called `Varnish`.

In another terminal port forward your Varnish service:

```bash
$ kubectl port-forward svc/varnish-test-varnish 8080:80
```

Make some requests to see metrics visualized in your Grafana dashboard.

The chart is not a complete solution and intended to be modified to the end user needs.