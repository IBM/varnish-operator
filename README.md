# Varnish Operator

[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/5895/badge)](https://bestpractices.coreinfrastructure.org/projects/5895)

### Project status: alpha
The project is in development and breaking changes can be introduced.

The purpose of the project is to provide a convenient way to deploy and manage Varnish instances in Kubernetes.

Kubernetes version `>=1.16.0` is supported.

Varnish version 6.5.1 is supported.

Full documentation can be found [here](https://ibm.github.io/varnish-operator/)

### Overview

Varnish operator manages Varnish clusters using a CustomResourceDefinition that defines a new Kind called `VarnishCluster`. 

The operator manages the whole lifecycle of the cluster: creating, deleting and keeping the cluster configuration up to date. The operator is responsible for building the VCL configuration using templates defined by the users and keeping the configuration up to date when relevant events occur (backend pod failure, scaling of the deployment, VCL configuration change).

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

### Further reading

* [QuickStart](https://ibm.github.io/varnish-operator/quick-start.html)
* [Contributing](https://ibm.github.io/varnish-operator/development.html)
