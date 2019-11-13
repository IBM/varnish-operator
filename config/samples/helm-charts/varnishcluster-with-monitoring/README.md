# VarnishCluster with monitoring

This chart's purpose is to show an example of how a `VarnishCluster` can be packaged into a Helm chart with the monitoring stack included. Feel free to copy the chart and modify it to your needs.

The monitoring stack is built using [Prometheus operator](https://github.com/helm/charts/tree/master/stable/prometheus-operator) and assumes it is installed already in the cluster. The Varnish operator should [be installed](https://pages.github.ibm.com/TheWeatherCompany/icm-varnish-k8s-operator/installation.html) as well.

The chart does the following:

* Creates a simple `VarnishCluster`
* Installs and configures Prometheus to scrape Varnish metrics
* Installs and configures Grafana to visualize scraped Prometheus metrics. A sample dashboard is included and available in Grafana after installation.

The example is not a complete solution and does not include PV/PVC creation, secure credentials setup, ingresses, etc. Those parts are left for the user to setup.

## Installation

You will need a backend to actually test if Varnish is working correctly. Use your existing or create one. This example will use a simple nginx server as the backend:

```bash
$ kubectl create deployment nginx --image nginx
``` 

Make sure you have the [imagePullSecret created](https://pages.github.ibm.com/TheWeatherCompany/icm-docs/managed-kubernetes/container-registry.html#pulling-an-image-in-kubernetes) to be able to pull Varnish images.

Clone the repo and install the chart using the local path to the chart:

```bash
$ git clone git@github.ibm.com:TheWeatherCompany/icm-varnish-k8s-operator.git
$ cd icm-varnish-k8s-operator
$ helm install --name varnish-test config/samples/helm-charts/varnishcluster-with-monitoring --set varnish.imagePullSecret=docker-reg-secret --set varnish.backendsSelector.app=nginx --set varnish.backendsPort=80
```

Note that we've specified the selector for our backends (`--set backendsSelector.app=nginx`) and the port they are listening on (`--set backendsPort=80`)

You should see your `VarnishCluster`, `Prometheus` and `Grafana` pods starting/running:

```bash
$ kubectl get pods                   
NAME                                                      READY   STATUS    RESTARTS   AGE
nginx-5c7588df-8wsr6                                      1/1     Running   0          5m16s
prometheus-varnish-test-prometheus-0                      3/3     Running   0          3m7s
varnish-test-grafana-9f584598d-89smp                      2/2     Running   0          3m8s
varnish-test-varnish-varnish-1                1/1     Running   0          3m7s
```

and the corresponding services created:

```bash
$ kubectl get svc                                            
NAME                            TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)           AGE
prometheus-operated             ClusterIP   None             <none>        9090/TCP          4m25s
varnish-test-grafana            ClusterIP   10.111.214.241   <none>        80/TCP            4m26s
varnish-test-prometheus         ClusterIP   10.97.244.131    <none>        9090/TCP          4m26s
varnish-test-varnish            ClusterIP   10.100.142.240   <none>        9131/TCP,80/TCP   4m25s
varnish-test-varnish-no-cache   ClusterIP   10.98.10.22      <none>        80/TCP            4m25s
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

Make a request to see metrics appear in your dashboard:

```bash
$ curl localhost:8080                       
  <!DOCTYPE html>
  <html>
  <head>
  <title>Welcome to nginx!</title>
  <style>
      body {
          width: 35em;
          margin: 0 auto;
          font-family: Tahoma, Verdana, Arial, sans-serif;
      }
  </style>
  </head>
  <body>
  <h1>Welcome to nginx!</h1>
  <p>If you see this page, the nginx web server is successfully installed and
  working. Further configuration is required.</p>
  
  <p>For online documentation and support please refer to
  <a href="http://nginx.org/">nginx.org</a>.<br/>
  Commercial support is available at
  <a href="http://nginx.com/">nginx.com</a>.</p>
  
  <p><em>Thank you for using nginx.</em></p>
  </body>
  </html>
```

## Uninstall

To uninstall the chart simply delete the chart:

```bash
helm delete --purge varnish-test
```