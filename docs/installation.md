# Varnish Operator 

## Installation

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

Use the image pull secret created in the previous step to install the operator:

```bash
$ helm install --name varnish-operator --namespace varnish-operator --set container.imagePullSecret=container-reg-secret icm/varnish-operator
```                                                                                                                        

You should see your operator pod up and running:

```bash
$ kubectl get pods --namespace varnish-operator
NAME                 READY   STATUS              RESTARTS   AGE
varnish-operator-0   1/1     Running             0          40s
```

Also, your cluster should have a new resource - `varnishservice` (with aliases `varnishservices` and `vs`). Now you're ready to create `VarnishService` resources.

To tolerate pod failures, you can run multiple replicas of the operator by specifying the `replicas` field to more than `1` in Helm value overrides. Only one of the  replicas - the leader - will handle the work while others monitor if the leader is healthy. If the leader fails, one of the backup replicas will start to handle the work. That also means that if you need to check operator logs for debugging purposes you need to find the leader pod and check its logs.

See [Operator Configuration](operator-configuration.md) section for more options to configure.

## Operator Update

Since the operator is packaged into a Helm chart, the update is done by a simple `helm upgrade` command.

Note that when the operator version updates, the Varnish images get updated as well which causes a restart of all Varnish pods controlled by the operator. Since Varnish is an in-memory cache it means that **all cached data will be lost** so plan your upgrade procedure accordingly.

## Uninstalling the Operator

Uninstallation of the chart is also done by Helm.
It deletes all created resources, including the CRD which will cause **deletion of all your VarnishService instances**.
