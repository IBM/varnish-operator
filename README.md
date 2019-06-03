# The VarnishService Kubernetes Operator

[![Build Status](https://wcp-twc-icmkube-jenkins.swg-devops.com/job/TheWeatherCompany%20ICM/job/icm-varnish-k8s-operator/job/master/badge/icon)](https://wcp-twc-icmkube-jenkins.swg-devops.com/job/TheWeatherCompany%20ICM/job/icm-varnish-k8s-operator/job/master/)

Documentation can be found [here](https://pages.github.ibm.com/TheWeatherCompany/icm-varnish-k8s-operator/).

VarnishService fills a space currently missing within Kubernetes on IBM Cloud: Varnish. IBM does not provide any managed Varnish instances, and Kubernetes does not have anything that works like Varnish does. Thus, this project aims to fill that space by providing a convenient way to deploy Varnish instances.

By default, deploying a Varnish directly as a Deployment into Kubernetes is not immediately useful because the VCL must have IP addresses for its backends. The only obvious way to get an IP address is via a Kubernetes Service, but that Service already acts as a load balancer to the Deployment it backs, which means undefined behavior from the Varnish perspective, and adds an extra network hop. Thus, trying to use Varnish in a regular deployment is unproductive.

Instead, the VarnishService operator manages the required infrastructure of your varnish, including the deployment, filling in the IP addresses of the pods for you. The operator is made up of 4 components:

**CustomResourceDefinition**: the actual "VarnishService" resource, that acts in the same way that a Service resource does, except with an added Varnish layer between the Service and the Deployment it backs. You would define a resource of Kind "VarnishService", and specify all the regular specs for a Service, plus some new fields that control how many Varnish instances you want, how much memory/cpu they get, and other relevant information for the Varnish cluster.

**Controller**: The controller is an application deployed into your cluster that knows how to react to the VarnishService CustomResource. Meaning, this application watches for new or changed VarnishServices and handles the actual underlying infrastructure. That means it must be running at all times in the cluster, although it lives in its own namespace away from your application.

**Varnish**: Of course, there will be a cluster of Varnishes. Currently, only Varnish version 6.1.1 is supported, but more versions will be supported in the future.

**K-Watcher**: K-watcher is a sidecar to every Varnish node that monitors the backends and updates the Varnish VCL whenever it notices a change to the pod IP addresses.

## Kubernetes Version Requirement

This operator assumes that the `/status` and `/scale` subresources are enabled for Custom Resources, which means that you must have enabled this alpha feature for Kubernetes v1.10 (impossible on IBM Kubernetes Service) or are using at least v1.11, where it is enabled by default.

## Installation

The VarnishService Operator is packaged as a [Helm Chart](https://helm.sh/), hosted on [Artifactory](https://na.artifactory.swg-devops.com). To get access to this Artifactory, you must be a user on the ICM Core Engineering Bluemix account 1638245, and specifically a Blue Group that has access to the Artifactory resources. 

### Getting Helm Access

After you are a user on the correct Blue Group, you must generate an API key within [Artifactory](https://na.artifactory.swg-devops.com) for Helm to use. You can generate an API key on your profile page, found in the upper-right of the home page. Using that generated API Key, you can log in to Helm using [these instructions](https://www.jfrog.com/confluence/display/RTF/Helm+Chart+Repositories), where the username is your email and the password is your API key. Specifically, that will look like:

```sh
helm repo add wcp-icm-helm-virtual https://na.artifactory.swg-devops.com/artifactory/wcp-icm-helm-virtual --username=<your-email> --password=<encrypted-password>
helm repo update
```

#### Remote Repo

Alternatively, you can configure the ICM Core Engineering Artifactory repo as a remote repo on your Artifactory repo. This has the advantage that members of the team will only need to configure access to a single repo (your team's). In order to configure a remote repo, you will need to talk directly to the TaaS team to have them set this up.

### Getting Container Registry Access

As part of the install, you will also need access to the Container Registry in order to pull the Docker images associated with the Helm charts. Instructions for this can be found in [the icm-docs](https://pages.github.ibm.com/TheWeatherCompany/icm-docs/managed-kubernetes/container-registry.html#pulling-an-image-in-kubernetes)

When installing, the Helm charts will create the operator (and thus look for the container-registry secret) in the `varnish-operator-system` namespace by default, although this can be overridden in a `values.yaml`.

### Configuring The Operator

The operator has options to customize the installation into your cluster, exposed as values in the Helm `values.yaml` file. [See the default `values.yaml` annotated with descriptions of each field](/varnish-operator/values.yaml) to see what can be customized when deploying this operator. You only need to specify values that are different from those contained in the `values.yaml`.

### Installing The Operator

Once a Namespace has been created with a docker registry secret and an appropriate `values.yaml` has been assembled, install the operator using

```sh
helm upgrade --install <name-of-release> wcp-icm-helm-virtual/varnish-operator --version <latest-version> --wait --namespace <namespace-with-registry-token>
```

Note that

* `<name-of-release>` can be any name and has the same meaning as `<name>` for `helm install --name <name>`. For consistency, you might consider using `varnish-operator`
* `<namespace-with-registry-token>` must match `namespace` in the `values.yaml` file.

## Usage

Once the operator is installed, your cluster should have a new resource, `varnishservice` (with aliases `varnishservices` and `vs`). From this point, you can create a yaml file with the `VarnishService` Kind.

### Configuring Access

Since the VarnishService requires pulling images from the same private repository as the Operator, the same docker registry key must exist in the target namespace for the VarnishService. Thus, add a secret with the docker registry token to that namespace before creating the resource.

### Configuring The VarnishService Resource

VarnishService has [an example yaml file annotated with descriptions of each field](/config/samples/icm_v1alpha1_varnishservice.yaml) To see what can be customized for the VarnishService. Copy this file and customize it to your needs.

### Preparing VCL Code

There are 2 fields relevant to configuring the VarnishService for VCL code, in `.spec.vclConfigMap`:

* **name**: This is a REQUIRED field, and tells the VarnishService the name of the ConfigMap that contains/will contain the VCL files
* **entrypointFile**: The name of the file that acts as the entrypoint for Varnish. This is the name of the file that will be passed to the Varnish executable.
  * If `entrypointFile` is templated (ends in `.tmpl`), exclude the `.tmpl` extension. eg: if ConfigMap has file `mytemplatedfile.vcl.tmpl`, set `entrypointFile: mytemplatedfile.vcl`

If a ConfigMap of name `.spec.vclConfigMap.name` does not exist on VarnishService creation, the operator will create one and populate it with a default `backends.vcl.tmpl` and `default.vcl`. Their behavior are as follows:

* [`backends.vcl.tmpl`](/config/vcl/backends.vcl.tmpl): collect all backends into a single director and round-robin between them
* [`default.vcl`](/config/vcl/default.vcl):
  * respond to `GET /heartbeat` checks with a 200
  * respond to `GET /liveness` checks with a 200 or 503, depending on healthy backends
  * respond to all other requests normally, caching all non-404 responses
  * hash request based on url
  * add `X-Varnish-Cache` header to response with "HIT" or "MISS" value, based on presence in cache

If you would like to use the default `backends.vcl.tmpl`, but a custom `default.vcl`, the easiest way is to create the VarnishService without the ConfigMap, let the operator create the ConfigMap for you, and then modify the contents of the ConfigMap after creation. Alternatively, just copy the content as linked above.

### Writing a Templated VCL File

The template file is a regular vcl file, with the addition of [Go templates](https://golang.org/pkg/text/template). This is because there is no way to know at startup what the IP addresses of the backends will be, so they must be injected at runtime. Not to mention, they can change over time if the backends get rescheduled by Kubernetes. These are the available fields in the template:

* .Backends - `[]PodInfo`: array of backends
  * .IP - `string`: ip address of a backend
  * .NodeLabels - `map[string]string`: labels of the node on which the backend is deployed. This is primarily for configuration of multi-zone clusters
  * .PodName - `string`: name of pod representing a backend
* .TargetPort - `int`: port that is exposed on the backends
* .VarnishNodes - `[]PodInfo`: array of varnish nodes, for configuration of shard director (if using round robin director, you can ignore)
  * .IP - `string`: ip address of a varnish node
  * .NodeLabels - `map[string]string`: labels of the k8s node on which a varnish node is deployed. This is primarily for configuration of multi-zone clusters
  * .PodName - `string`: name of pod representing a varnish node
* .VarnishPort - `int`: port that is exposed on the varnish nodes (if using round robin director, you can ignore)

For example, to loop over the backends and create vcl `backend`s for each:

```vcl
{{ range .Backends }}
backend {{ .PodName }} {
  .host = "{{ .IP }}";
  .port = "{{ $.TargetPort }}";
}
{{ end }}
```

This loops over `.Backends`, names each backend `.PodName`, sets `.host` to `.IP`, and then sets port to the universal `$.TargetPort`.

For the full example of using the templates, see the [`backends.vcl.tmpl` file](/config/vcl/backends.vcl.tmpl).

### Using User Defined VCL Code Versions

VCL related status information is available at field `.status.vcl`. 

The current VCl version can be found at `.status.vcl.configMapVersion`. It matches the resource version of the config map that contains the VCL code. 

To tag your own versions, an annotation `VCLVersion` on the ConfigMap can be used.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  annotations:
    VCLVersion: v1.0 # <-- set by the user
  resourceVersion: "292181"
  ...
data:
    ...
```

After setting the annotation, that version can be seen at `.status.vcl.version`. This field is optional, and will only appear if the annotation is set.

```yaml
apiVersion: icm.ibm.com/v1alpha1
kind: VarnishService
metadata:
    ...
status:
  vcl:
    version: v1.0 # <-- reflects the `VCLVersion` annotation in the config map
    configMapVersion: "292181" # <-- reflects the config map resource version
  ...
```

After the VCL in the ConfigMap has been changed, the associated status fields will be immediately updated to reflect the latest version. However that does not guarantee that Varnish pods run the latest VCL configuration. It needs time to reload and if there is a problem, such as a syntax error in the VCL, may never load.
 
For better observability about currently running VCL versions, see `.status.vcl.availability`, which indicates how many pods have the latest version and how many of them are outdated. 

```yaml
apiVersion: icm.ibm.com/v1alpha1
kind: VarnishService
metadata:
    ...
status:
  vcl:
    configMapVersion: "292181"
    version: v1.0
    availability: 1 latest / 0 outdated # <-- all pods have the latest VCL version
  deployment:
    availableReplicas: 1
    ...
```

To check which pods have outdated versions, simply check their annotations. The annotation `configMapVersion` on the Varnish pod will indicate the latest version of the ConfigMap used. If it's not the same as in the VarnishService status it's likely that there's an issue.

Example of detecting a pod that failed to reload:

```bash
# get the latest version
> kubectl get varnishservice -n varnish-ns my-varnish -o=custom-columns=NAME:.metadata.name,CONFIG_MAP_VERSION:.status.vcl.configMapVersion
NAME        CONFIG_MAP_VERSION
my-varnish  292181
# figure out which pods doesn't have that latest version
> kubectl get pods -n varnish-ns -o=custom-columns=NAME:.metadata.name,CONFIG_MAP_VERSION:.metadata.annotations.configMapVersion
NAME                                            CONFIG_MAP_VERSION
my-varnish-varnish-deployment-545f475b58-7xn9k  292181
my-varnish-varnish-deployment-545f475b58-jc5vg  292181
my-varnish-varnish-deployment-545f475b58-nqqd2  351231 #outdated VCL
# check logs for that pod with outdated VCL
> kubectl logs -n my-varnish my-varnish-varnish-deployment-545f475b58-nqqd2 
2018-12-21T17:03:07.917Z	INFO	controller/controller.go:124	Rewriting file	{"path": "/etc/varnish/backends.vcl"}
2018-12-21T17:03:17.904Z	ERROR	controller/controller.go:157	exit status 1
/go/src/icm-varnish/k-watcher/pkg/controller/controller_varnish.go:50: Message from VCC-compiler:
Expected one of
	'acl', 'sub', 'backend', 'probe', 'import', 'vcl',  or 'default'
Found: 'dsafdf' at
('/etc/varnish/backends.vcl' Line 4 Pos 2)
 dsafdf
-######

Running VCC-compiler failed, exited with 2
Command failed with error code 106
VCL compilation failed
No VCL named v304255 known.
Command failed with error code 106

/go/src/icm-varnish/k-watcher/vendor/sigs.k8s.io/controller-runtime/pkg/internal/controller/controller.go:207: 
icm-varnish/k-watcher/pkg/logger.WrappedError
	/go/src/icm-varnish/k-watcher/pkg/logger/logger.go:49
ic
```

As the logs indicate, the issue here is the invalid VCL syntax.

### Creating a VarnishService Resource

Once the VarnishService resource yaml is ready, simply `kubectl apply -f <varnish-service>.yaml` to create the resource. Once complete, you should see:

* a deployment with the name `<varnish-service-name>-deployment`. This is the Varnish cluster, and should have inherited everything under the `deployment` part of the spec.
* 2 services, one `<varnish-service-name>` and one `<varnish-service-name>-no-cache`. As is implied by the names, using `<varnish-service-name>` will act as the service configured under `.spec.service`, and will direct to Varnish before hitting the underlying deployment, while `<varnish-service-name>-no-cache` will target the underlying deployment directly, with no Varnish caching. `<varnish-service-name>` will have inherited everything under the `service` part of the spec, other than its `selector`, which will be redirected to the Varnish deployment.
* A ConfigMap with VCL in it (either user-created, before running `kubectl apply -f <varnish-service>.yaml`, or generated by operator)
* A role/rolebinding/clusterrole/clusterrolebinding/serviceAccount combination to give the Varnish deployment the ability to access necessary resources.
* If configured, a PodDisruptionBudget as specced

### Updating a VarnishService Resource

Just as with any other Kubernetes resource, using `kubectl apply`, `kubectl patch`, or `kubectl replace` will all update the VarnishService appropriately. The operator will handle how that update propagates to its dependent resources. Conversely, trying to modify any of those dependent resources (Deployment, Services, Roles/Rolebindings, etc) will cause the operator to revert those changes, in the same way a Deployment does for its Pods. The only exception to this is the ConfigMap, the contents of which you can and should modify, since that is the VCL used to run the Varnish Pods.

### Deleting a VarnishService Resource

Simply calling `kubectl delete` on the VarnishResource will recursively delete all dependent resources, so that is the only action you need to take. This includes a user-generated ConfigMap, as the VarnishService will take ownership of that ConfigMap after creation. Deleting any of the dependent resources will trigger the operator to recreate that resource, in the same way that deleting the Pod of a Deployment will trigger the recreation of that Pod.

### Checking Status of a VarnishService Resource

The VarnishService keeps track of its current status as events occur in the system. This can be seen through the `Status` field, visible from `kubectl describe vs <your-varnishservice>`.

### Prove that Varnish is Working

In order to show that Varnish is properly configured, it is common to add a header to the response indicating whether Varnish responded from cache or from origin. In the `default.vcl` provided out of the box, that header is `X-Varnish-Cache: HIT` and `X-Varnish-Cache: MISS`. With such a header prepared, make a request against the service (ie, name matches your `<varnish-service>`) and look at its headers for the varnish-added header.

To make such a request, you can either

* expose that service through ingress, so that it is accessible to the outside world
* exec onto a pod specifically to make a curl request against the service. For instance, running

    ```sh
    kubectl run curlbox --image=radial/busyboxplus:curl --restart=Never -it -- sh
    ```

    will log you into a pod in the cluster that has curl and can access the service by name
  * After exiting from the `curlbox` container, run `kubectl delete pod curlbox` to clean it up 

## Keeping Varnish Stable

Kubernetes is built on the premise that its runnable environments are ephemeral, meaning they can be created or deleted at will, with little to no effect on the overall system. In the case of Varnish, which is purely an in-memory caching layer, deleting and creating instances all the time would cause the cache to perform very poorly. Thus, there is a need to keep Varnish stable, ie tell Kubernetes that these particular runnable environments should _not_ be treated as ephemeral.

Kubernetes does not provide this functionality out of the box, but you can trick it into approximating this behavior, and that is through the concepts of guaranteed resources and affinities.

### Guaranteed Resources

The way that Kubernetes manages deployed pods on nodes is through monitoring the resources that a pod is using. Specifically, it uses the `limits` and `requests` values for `cpu` and `memory` to determine how much resources to give a pod, and when it might be OK to reschedule a pod somewhere else (namely, if a node is running out of resources and some pods are using more resources than requested). For a detailed breakdown of what `limits` and `requests` mean, [see the Kubernetes documentation on QoS](https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/). In QoS parlance, you want the Varnish nodes to be a "Guaranteed" QoS. In short, you want to always set the `limits` and `requests` fields, and you want `limits` and `requests` to be identical.

### Affinities

[Kubernetes allows control](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#inter-pod-affinity-and-anti-affinity-beta-feature) on where pods get deployed based on labels associated with pods and nodes. For instance, you can configure pods of the same deployment to repel each other, meaning new pods entering the deployment will try to avoid nodes that already have a pod of that type. That way, you if any one node goes down, it will only take a single pod with it. Likewise, you can configure pods to be attracted to each other, for colocation that could decrease latency between pods. Note that reading through the above linked documentation is valuable, as it goes into limitations to affinities, as well as deeply explains how they work and when to use them.

For the purposes of this Varnish deployment, you will most likely want to configure a [pod anti-affinity](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#never-co-located-in-the-same-node) so that each pod of the varnish deployment is on a different node. An example of what that might look like is in the [example annotated yaml file](/config/samples/icm_v1alpha1_varnishservice.yaml) under `spec.deployment.affinity`.

### Running Varnish pods on separate IKS worker pools

This example shows how to create an IKS worker pool and make Varnish pods run strictly on its workers, one per node.

References:
 * [How to create IKS clusters and worker pools.](https://console.bluemix.net/docs/containers/cs_clusters.html#clusters)
 * [Taints and Tolerations](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/)
 * [Affinity and anti-affinity](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity)
 
1. Create a worker pool in your cluster assuming you already have a cluster called `test-cluster`

    ```bash
    $ #Find out the available zones to your cluster
    $ ibmcloud ks cluster-get --cluster test-cluster | grep "Worker Zones" # Get the 
    Worker Zones:           dal10
    $ #Find out what machine type are available in your zone  
    $ ibmcloud ks machine-types --zone dal10
    OK
    Name                      Cores   Memory   Network Speed   OS             Server Type   Storage   Secondary Storage   Trustable   
    u2c.2x4                   2       4GB      1000Mbps        UBUNTU_16_64   virtual       25GB      100GB               false   
    ms2c.4x32.1.9tb.ssd       4       32GB     10000Mbps       UBUNTU_16_64   physical      2000GB    960GB               false   
    ms2c.16x64.1.9tb.ssd      16      64GB     10000Mbps       UBUNTU_16_64   physical      2000GB    960GB               true   
    ms2c.28x256.3.8tb.ssd     28      256GB    10000Mbps       UBUNTU_16_64   physical      2000GB    1920GB              true   
       ...
    $ #Create a worker pool. 
    $ ibmcloud ks worker-pool-create --name varnish-worker-pool --cluster test-cluster --machine-type u2c.2x4 --size-per-zone 2 --hardware shared
    OK 
    $ #Verify your worker pool is created
    $ ibmcloud ks worker-pools --cluster test-cluster
    Name                  ID                                         Machine Type          Workers   
    default               91ed9433e7bf4dc7b8348ae1022f9f27-89d7d12   b2c.16x64.encrypted   3   
    varnish-worker-pool   91ed9433e7bf4dc7b8348ae1022f9f27-c5b13da   u2c.2x4.encrypted     2   
    OK
    $ #Add your zone to your worker pool. First, find out your VLAN IDs
    $ ibmcloud ks vlans --zone dal10
    OK
    ID        Name   Number   Type      Router         Supports Virtual Workers   
    2315193          1690     private   bcr02a.dal10   true   
    2315191          1425     public    fcr02a.dal10   true
    $ #Use the VLAN IDs above to add your zone to the worker pool
    $ ibmcloud ks zone-add --zone dal10 --cluster test-cluster --worker-pools varnish-worker-pool --private-vlan 2315193 --public-vlan 2315191
    OK
    $ #Verify that worker nodes provision in the zone that you've added
    $ ibmcloud ks workers --cluster test-cluster --worker-pool varnish-worker-pool
    OK
    ID                                                  Public IP   Private IP   Machine Type        State               Status                          Zone    Version   
    kube-dal10-cr91ed9433e7bf4dc7b8348ae1022f9f27-w58   -           -            u2c.2x4.encrypted   provision_pending   Preparing to provision worker   dal10   1.11.7_1543   
    kube-dal10-cr91ed9433e7bf4dc7b8348ae1022f9f27-w59   -           -            u2c.2x4.encrypted   provision_pending   -                               dal10   1.11.7_1543   
    ```
    
    Wait until your worker pool nodes change their state to `normal` and status to `Ready`.
    
    ```bash
    $ ibmcloud ks workers --cluster test-cluster --worker-pool varnish-worker-pool
    OK
    ID                                                  Public IP       Private IP      Machine Type        State    Status   Zone    Version   
    kube-dal10-cr91ed9433e7bf4dc7b8348ae1022f9f27-w58   169.61.218.68   10.94.177.179   u2c.2x4.encrypted   normal   Ready    dal10   1.11.7_1543   
    kube-dal10-cr91ed9433e7bf4dc7b8348ae1022f9f27-w59   169.61.218.94   10.94.177.180   u2c.2x4.encrypted   normal   Ready    dal10   1.11.7_1543
    ```
    
1. Taint created nodes to repel pods that don't have required toleration. 

    ```bash
    $ #Setup kubectl
    $ ibmcloud ks cluster-config --cluster test-cluster 
    OK
    The configuration for test-cluster was downloaded successfully.
    
    Export environment variables to start using Kubernetes.
    
    export KUBECONFIG=/home/me/.bluemix/plugins/container-service/clusters/test-cluster/kube-config-dal10-test-cluster.yml
    
    $ export KUBECONFIG=/home/me/.bluemix/plugins/container-service/clusters/test-cluster/kube-config-dal10-test-cluster.yml
    $ #Find your nodes using kubectl. First get your worker pool ID and then use it to select your nodes
    $ ibmcloud ks worker-pools --cluster test-cluster 
    Name                  ID                                         Machine Type          Workers   
    default               91ed9433e7bf4dc7b8348ae1022f9f27-89d7d12   b2c.16x64.encrypted   3   
    varnish-worker-pool   91ed9433e7bf4dc7b8348ae1022f9f27-c5b13da   u2c.2x4.encrypted     2   
    $ kubectl get nodes --selector ibm-cloud.kubernetes.io/worker-pool-id=91ed9433e7bf4dc7b8348ae1022f9f27-c5b13da
    NAME            STATUS   ROLES    AGE   VERSION
    10.94.177.179   Ready    <none>   16m   v1.11.7+IKS
    10.94.177.180   Ready    <none>   15m   v1.11.7+IKS
    $ #Taint those nodes
    $ kubectl taint node 10.94.177.179 role=varnish:NoSchedule #Do not schedule here not Varnish pods
    node/10.94.177.179 tainted
    $ kubectl taint node 10.94.177.179 role=varnish:NoExecute #Evict not Varnish pods if they already running here
    node/10.94.177.179 tainted
    $ kubectl taint node 10.94.177.180 role=varnish:NoSchedule #Do not schedule here not Varnish pods
    node/10.94.177.180 tainted
    $ kubectl taint node 10.94.177.180 role=varnish:NoExecute #Evict not Varnish pods if they already running here
    node/10.94.177.180 tainted
    ```
    
    This prevents all pods from scheduling on that node unless you already have pods with matching toleration
    
1. Label the nodes for the ability to schedule your varnish pods only on that nodes. Those labels will be used in your VarnishService configuration later.

    ```bash
    $ kubectl label node 10.94.177.179 role=varnish-cache
    node/10.94.177.179 labeled
    $ kubectl label node 10.94.177.180 role=varnish-cache
    node/10.94.177.180 labeled 
    ```
1. Define your VarnishService spec with necessary affinity and toleration configuration

    4.1 Define pods anti-affinity to not co-locate replicas on a node.
    
    ```yaml
    metadata:
      labels:
        role: varnish-cache
    spec:
      deployment:
        affinity:
          podAntiAffinity:
            requiredDuringSchedulingIgnoredDuringExecution:
              - labelSelector:
                  matchExpressions:
                    - key: role
                      operator: In
                      values:
                        - varnish-cache
                topologyKey: "kubernetes.io/hostname"
    ```
    That will make sure that two varnish pods doesn't get scheduled on one node. Kubernetes makes the decision based on labels we've set in the spec
    
    4.2 Define pods node affinity
    
    ```yaml
    spec:
      deployment:
        affinity:
          nodeAffinity:
            requiredDuringSchedulingIgnoredDuringExecution:
              nodeSelectorTerms:
              - matchExpressions:
                - key: role
                  operator: In
                  values:
                    - varnish-cache
    ```
    That will make kubernetes schedule varnish pods only on our worker pool nodes. The labels used here are the ones we've assigned to the node in step 3
    
    4.3 Define pods tolerations
    
    ```yaml
    spec:
      deployment:
        tolerations:
          - key: "role"
            operator: "Equal"
            value: "varnish"
            effect: "NoSchedule"
          - key: "role"
            operator: "Equal"
            value: "varnish"
            effect: "NoExecute"
    ```
    In step 2 we made our node repel all pods that don't have specific tolerations. Here we added those tolerations to be eligible for scheduling on those nodes. The values are the ones we used when tainted our nodes in step 2. 
    
5. Apply your configuration.

    This step assumes you have varnish operator [installed](#installation) and the namespace has the necessary secret [installed](#configuring-access).
    
    Complete VarnishService configuration example:
    
    ```yaml
    apiVersion: icm.ibm.com/v1alpha1
    kind: VarnishService
    metadata:
      labels:
        role: varnish-cache
      name: varnish-in-worker-pool
      namespace: varnish-ns
    spec:
      vclConfigMap:
        name: varnish-worker-pool-files
        backendsFile: backends.vcl
        defaultFile: default.vcl
      deployment:
        replicas: 2
        container:
          resources:
            limits:
              cpu: 100m
              memory: 256Mi
            requests:
              cpu: 100m
              memory: 256Mi
          readinessProbe:
            exec:
              command: [/usr/bin/varnishadm, ping]
          imagePullSecret: docker-reg-secret
        affinity:
          podAntiAffinity:
            requiredDuringSchedulingIgnoredDuringExecution:
              - labelSelector:
                  matchExpressions:
                    - key: role
                      operator: In
                      values:
                        - varnish-cache
                topologyKey: "kubernetes.io/hostname"
          nodeAffinity:
            requiredDuringSchedulingIgnoredDuringExecution:
              nodeSelectorTerms:
              - matchExpressions:
                - key: role
                  operator: In
                  values:
                    - varnish-cache
        tolerations:
          - key: "role"
            operator: "Equal"
            value: "varnish-cache"
            effect: "NoSchedule"
          - key: "role"
            operator: "Equal"
            value: "varnish-cache"
            effect: "NoExecute"
      service:
        selector:
          app: HttPerf
        varnishPort:
          name: varnish
          port: 2035
          targetPort: 8080
        varnishExporterPort:
          name: varnishexporter
          port: 9131
    ```
    Apply your configuration:
    ```bash
    $ kubectl apply -f varnish-in-worker-pool.yaml
    varnishservice.icm.ibm.com/varnish-in-worker-pool created
    ```
    Here the operator will create all pods with specified configuration
6. See your pods being scheduled strictly on your worker pool and spread across different nodes.
    ```bash
    $ kubectl get pods -n varnish-ns -o wide --selector role=varnish-cache
    NAME                                                         READY   STATUS    RESTARTS   AGE   IP               NODE            NOMINATED NODE
    varnish-in-worker-pool-varnish-deployment-78c9b6f5bf-kqg72   1/1     Running   0          6m    172.30.244.65    10.94.177.179   <none>
    varnish-in-worker-pool-varnish-deployment-78c9b6f5bf-pqtzv   1/1     Running   0          6m    172.30.136.129   10.94.177.180   <none>

    ```
    Check the `NODE` column. The value will be different for each pod.
    
    Note that you won't be able to run more pods than you have nodes. The anti-affinity rule will not allow two pods being co-located on one node.
    This behaviour can be changed by using an anti-affinity type called `preferredDuringSchedulingIgnoredDuringExecution`: 
    
    ```yaml
    spec:
      deployment:
        affinity:
          podAntiAffinity:
            preferredDuringSchedulingIgnoredDuringExecution:
            - weight: 1
              podAffinityTerm:
                labelSelector:
                  matchExpressions:
                  - key: role
                    operator: In
                    values:
                    - varnish-cache
                topologyKey: "kubernetes.io/hostname"
    ```
     It will still ask Kubernetes to spread pods onto different nodes but also allow to co-locate them if there are more pods than nodes.
     
    Also keep in mind that in such configuration the pods can be scheduled to your worker pool only. If the worker pool is deleted the pods will hang in `Pending` state until new nodes with the same configuration are added to the cluster.
