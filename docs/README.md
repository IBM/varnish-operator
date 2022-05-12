# Introduction

Varnish operator creates, configures and manages Varnish clusters. It generates the VCL configuration using user defined templates and keeps it up to date as the cluster changes.

## Features

 * [x] Basic install
 * [x] Full lifecycle support (create/update/delete)
 * [x] Automatic VCL configuration updates (using user defined templates)
 * [x] Prometheus metrics support
 * [x] Scaling
 * [x] Configurable update strategy
 * [x] Persistence (for [file storage backend](https://varnish-cache.org/docs/trunk/users-guide/storage-backends.html#file) support)
 * [ ] Multiple Varnish versions support
 * [ ] Autoscaling

### Overview

The operator works based on a CustomResourceDefinition that manages the Varnish cluster. It defines a new kind called `VarnishCluster` that describes the desired state of your Varnish instances.

Example of a simple `VarnishCluster`:

```yaml
apiVersion: caching.ibm.com/v1alpha1
kind: VarnishCluster
metadata:
  labels:
    operator: varnish
  name: varnishcluster-sample
  namespace: varnish-ns
spec:
  vcl:
    configMapName: vcl-files
    entrypointFileName: entrypoint.vcl
  replicas: 3
  backend:  
    selector:
      app: nginx
  service:
    port: 80
```

See the [VarnishCluster configuration section](varnish-cluster-configuration.md) for more details about the `VarnishCluster` spec.

### VCL configuration

The VCL configuration is generated using templated VCL files stored in a config map. [Go templates](https://golang.org/pkg/text/template/) are used as the template engine and can be used to generate Varnish backend definitions and to build your directors.

See the [VCL files configuration](vcl-configuration.md) section for more details.

### Incentive for building an operator

* Ability to manage Varnish instances in a Kubernetes native way 
* Deploying Varnish directly as a Deployment into Kubernetes is not immediately useful because the VCL has to know about the backend hosts. Those host names (or IP addresses) need to be stable in order to keep the VCL valid, but that's not possible due to dynamic nature of Kubernetes. The only obvious way to get a stable hostname (IP address) is via a Kubernetes Service, but that Service already acts as a load balancer to the Deployment it backs, which means undefined behavior from the Varnish perspective, and adds an extra network hop. Thus, trying to use Varnish in a regular deployment is unproductive.
* Support of different directors for backends. If you expose your backends to Varnish as a Kubernetes service you can have only round-robin load balancing. Since the operator works with backends at the pod level, you can use different directors supported by Varnish ([random](https://varnish-cache.org/docs/5.1/reference/vmod_directors.generated.html#obj-random), [fallback](https://varnish-cache.org/docs/5.1/reference/vmod_directors.generated.html#obj-fallback))
* You can't build a sharded Varnish cluster due to the dynamic nature of Kubernetes and the requirement to know about each Varnish peer (pod) in order to build the VCL.

### Further reading

* [Quickstart](quick-start.md)
* [VarnishCluster configuration](varnish-cluster-configuration.md)
* [Varnish operator configuration](operator-configuration.md)
* [VCL files configuration](vcl-configuration.md)
* [Contribution](development.md)
