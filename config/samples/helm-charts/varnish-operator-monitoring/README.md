# Monitoring stack for Varnish Operator

This chart's purpose is to show an example of a monitoring stack that monitors the Varnish Operator. Feel free to copy the chart and modify it to your needs.

The monitoring stack is built using the [Prometheus operator](https://github.com/helm/charts/tree/master/stable/prometheus-operator) and assumes it is installed already in the cluster.

The chart does the following:

* Installs and configures Prometheus to scrape Varnish Operator metrics
* Installs and configures Grafana to visualize scraped Prometheus metrics. A sample dashboard is included and available in Grafana after installation.

The example is not a complete solution and does not include PV/PVC creation, secure credentials setup, ingresses, etc. Those parts are left for the user to setup. 

## Installation

[Install the operator](https://pages.github.ibm.com/TheWeatherCompany/icm-varnish-k8s-operator/installation.html) if you haven't yet.

Clone the repo and install the chart using the local path to the chart:

```bash
$ git clone git@github.ibm.com:TheWeatherCompany/icm-varnish-k8s-operator.git
$ cd icm-varnish-k8s-operator
$ helm install --name varnish-operator-monitoring config/samples/helm-charts/varnish-operator-monitoring
```

Make sure to install the chart in the same namespaces as your operator.

You should see your `Prometheus` and `Grafana` pods starting/running:

```bash
$ kubectl get pods                                      
NAME                                                     READY   STATUS    RESTARTS   AGE
prometheus-varnish-operator-monitoring-prometheus-0      3/3     Running   1          37s
varnish-operator-0                                       1/1     Running   3          2d21h
varnish-operator-monitoring-grafana-6f7ff7f4f9-2pjpj     2/2     Running   0          37s
```

and the corresponding services created:

```bash
$ kubectl get svc                                            
NAME                            TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)           AGE
varnish-operator                          ClusterIP   10.102.21.96     <none>        9131/TCP            2d21h
varnish-operator-monitoring-grafana       ClusterIP   10.100.254.40    <none>        80/TCP              89s
varnish-operator-monitoring-prometheus    ClusterIP   10.109.153.192   <none>        9090/TCP            89s
varnish-operator-webhook-service          ClusterIP   10.103.72.178    <none>        443/TCP             15d
```

Port forward the Grafana pod to see the dashboard:

```bash
$ kubectl port-forward pod/varnish-operator-monitoring-grafana-6f7ff7f4f9-2pjpj 3000
Forwarding from 127.0.0.1:3000 -> 3000
Forwarding from [::1]:3000 -> 3000
```

You can see your dashboard at `localhost:3000`. The login is `admin`, password is `pass`. You will find a dashboard called `Varnish Operator`.

## Uninstall

To uninstall the chart simply delete the chart:

```bash
helm delete --purge varnish-operator-monitoring
```