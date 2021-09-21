# Quick Start

### Prerequisites

* Kubernetes v1.16 or newer and `kubectl` configured to communicate with your cluster
* Helm

### Configure Helm repo 

```bash
$ helm repo add varnish-operator https://raw.githubusercontent.com/IBM/varnish-operator/main/helm-releases
$ helm repo update
```

### Install Varnish Operator

```bash
$ helm install --name varnish-operator --namespace varnish-operator varnish-operator/varnish-operator
```                                                                                                                        

You should see your operator pod up and running:

```bash
$ kubectl get pods --namespace varnish-operator
NAME                              READY   STATUS              RESTARTS   AGE
varnish-operator-fd96f48f-gn6mc   1/1     Running             0          40s
```

### Create a VarnishCluster

1. Create a namespace where the `VarnishCluster` with the backend will live.

    ```bash
    $ kubectl create ns varnish-cluster
    ```

1. Create a simple backend that will be cached by Varnish:

    ```bash
    $ kubectl create deployment nginx-backend --image nginx -n varnish-cluster
      deployment.apps/nginx-backend created
    $ kubectl get deployment -n varnish-cluster nginx-backend --show-labels #get pod labels, they will be used to identify your backend pods
      NAME            READY   UP-TO-DATE   AVAILABLE   AGE   LABELS
      nginx-backend   1/1     1            1           15s   app=nginx-backend 
    ```

1. Create your `VarnishCluster`:

    ```bash
    $ cat <<EOF | kubectl create -f -
    apiVersion: caching.ibm.com/v1alpha1
    kind: VarnishCluster
    metadata:
      name: varnishcluster-example
      namespace: varnish-cluster # the namespace we've created above
    spec:
      vcl:
        configMapName: vcl-config # name of the config map that will store your VCL files. Will be created if doesn't exist.
        entrypointFileName: entrypoint.vcl # main file used by Varnish to compile the VCL code.
      backend:
        port: 80
        selector:
          app: nginx-backend # labels that identify your backend pods
      service:
        port: 80 # Varnish pods will be exposed on that port
    EOF
 
    varnishcluster.ibm.com/varnishcluster-example created  
    ```

Step 2 and 3 can be executed in any order. If the backend doesn't exist on `VarnishCluster` creation, Varnish will still work but won't cache any backends responding with 504 response code. Once the backend pods are up and running, Varnish will automatically pick up the changes, reload the VCL and start serving and caching the requests.

## What You Should See

Once `VarnishCluster` is created, you should see:

* StatefulSet called `<varnish-cluster-name>`
* Service called `<varnish-cluster-name>` which uses Varnish for caching
* Service called `<varnish-cluster-name>-no-cache` which bypasses Varnish
* ConfigMap called `vcl-config` containing VCL files that Varnish is using
* Role/Rolebinding/ClusterRole/ClusterRoleBinding/ServiceAccount combination for permissions

You can check if all works by doing `kubectl port-forward` and checking the server response.

Port forward your service:

```bash
$ kubectl port-forward -n varnish-cluster service/varnishcluster-example 8080:80
Forwarding from 127.0.0.1:8080 -> 6081
Forwarding from [::1]:8080 -> 6081
...
```

In another terminal, use `curl` to make a request to Varnish (use `-i` flag to see response headers):

```bash
$ curl -i localhost:8080/
  HTTP/1.1 200 OK
  Server: nginx/1.17.0
  Date: Tue, 11 Jun 2019 13:55:16 GMT
  Content-Type: text/html
  Content-Length: 612
  Last-Modified: Tue, 21 May 2019 14:23:57 GMT
  ETag: "5ce409fd-264"
  X-Varnish: 32770 12
  Age: 41
  Via: 1.1 varnish (Varnish/6.0)
  grace: 
  X-Varnish-Cache: HIT
  Accept-Ranges: bytes
  Connection: keep-alive
  
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

You should see nginx's welcome page. 

`Server: nginx/1.17.0` header shows your backend response header and `Via: 1.1 varnish (Varnish/6.1.1)` indicates that the request has been passed through Varnish.

## What's next

* [See how to configure your VCL files](vcl-configuration.md)
* [Configure your `VarnishCluster`](varnish-cluster-configuration.md) (resource requests and limits, affinity rules, tolerations)
* [Adjust your Varnish Operator configs](operator-configuration.md)
