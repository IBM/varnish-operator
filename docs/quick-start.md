# Quick Start

### Prerequisites

* Kubernetes v1.11 or newer and `kubectl` configured to communicate with your cluster
* Helm

### Configure Helm repo

Helm charts [are hosted in private Artifactory](https://pages.github.ibm.com/TheWeatherCompany/icm-docs/helm/chart-repositories.html#using-artifactory-as-a-helm-chart-repository) so you need to configure repo access first.

1. Get access to [Artifactory](https://na.artifactory.swg-devops.com)
1. Generate an API Key on the Artifactory website in your profile settings (click on your email in the top right corner)
1. Add the repo and update your local list of charts: 

    ```bash
    $ helm repo add icm https://na.artifactory.swg-devops.com/artifactory/wcp-icm-helm-virtual --username=<your-email> --password=<api-key>
    $ helm repo update
    ```
    
### Configure image pull secret

Create a namespace for the operator:

```bash
$ kubectl create ns varnish-operator
```

Images are located in a private IBM cloud registry. You need to [create an image pull secret](https://pages.github.ibm.com/TheWeatherCompany/icm-docs/managed-kubernetes/container-registry.html#pulling-an-image-in-kubernetes) in your namespace to be able to pull images in the cluster.

```bash
$ kubectl create secret docker-registry container-reg-secret \
    --namespace varnish-operator \
    --docker-server us.icr.io \
    --docker-username <user-name> \
    --docker-password=<password> \
    --docker-email=<email>
```

### Install Varnish Operator

Use the image pull secret created in the previous step to install your operator:

```bash
$ helm install --name varnish-operator --namespace varnish-operator --set container.imagePullSecret=container-reg-secret icm/varnish-operator
```                                                                                                                        

You should see your operator pod up and running:

```bash
$ kubectl get pods --namespace varnish-operator
NAME                 READY   STATUS              RESTARTS   AGE
varnish-operator-0   1/1     Running             0          40s
```

### Create a VarnishService

1. Create a namespace where the VarnishService with the backend will live.

    ```bash
    $ kubectl create ns varnish-service
    ```

1. Create the same image pull secret there. It will be used to pull Varnish images.

    ```bash
    $ kubectl create secret docker-registry container-reg-secret \
        --namespace varnish-service \
        --docker-server us.icr.io \
        --docker-username <user-name> \
        --docker-password=<password> \
        --docker-email=<email>
    ```
1. Create a simple backend that will be cached by Varnish:

    ```bash
    $ kubectl create deployment nginx-backend --image nginx -n varnish-service
      deployment.apps/nginx-backend created
    $ kubectl get deployment -n varnish-service nginx-backend --show-labels #get pod labels, they will be used to identify your backend pods
      NAME            READY   UP-TO-DATE   AVAILABLE   AGE   LABELS
      nginx-backend   1/1     1            1           15s   app=nginx-backend 
    ```

1. Create your VarnishService:

    ```bash
    $ cat <<EOF | kubectl create -f -
    apiVersion: icm.ibm.com/v1alpha1
    kind: VarnishService
    metadata:
      name: varnishservice-example
      namespace: varnish-service # the namespace we've created above
    spec:
      vclConfigMap:
        name: vcl-config # name of the config map that will store your VCl files. Will be created if doesn't exist.
        entrypointFile: entrypoint.vcl # main file used by Varnish to compile the VCL code.
      statefulSet:
        replicas: 3 # run 3 replicas of Varnish
        container: 
          imagePullSecret: container-reg-secret # the image pull secret created above
      service:
        selector:
          app: nginx-backend # labels that identify your backend pods
        varnishPort:
          name: varnish
          port: 80 # Varnish pods will be exposed on that port 
          targetPort: 80 # the port our backend pods listen on. 80 for nginx.
        varnishExporterPort: # prometheus exporter metrics port
          name: varnishexporter
          port: 9131
    EOF
 
    varnishservice.icm.ibm.com/varnishservice-example created  
    ```

## What You Should See

Once `VarnishService` is created, you should see:

* StatefulSet called `<varnish-service-name>-statefulset`
* Service called `<varnish-service-name>` which uses Varnish for caching
* Service called `<varnish-service-name>-no-cache` which bypasses Varnish
* ConfigMap called `vcl-config` containing VCL files that Varnish is using
* Role/Rolebinding/ClusterRole/ClusterRoleBinding/ServiceAccount combination for permissions

You can check if all works by doing `kubectl port-forward` and checking the server response.

Port forward your service:

```bash
$ kubectl port-forward -n varnish-service service/varnishservice-example 8080:80
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

`Server: nginx/1.17.0` header shows your backend response header and `Via: 1.1 varnish (Varnish/6.0)` indicates that the request has been passed through Varnish.

## What's next

* [See how to configure your VCL files](vcl-configuration.md)
* [Configure your `VarnishService`](varnish-service-configuration.md) (resource requests and limits, affinity rules, tolerations)
* [Adjust your Varnish Operator configs](operator-configuration.md)
