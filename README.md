# The VarnishService Kubernetes Operator

#### Build Status
[![Build Status](https://wcp-twc-icmkube-jenkins.swg-devops.com/job/TheWeatherCompany%20ICM/job/icm-varnish-k8s-operator/job/master/badge/icon)](https://wcp-twc-icmkube-jenkins.swg-devops.com/job/TheWeatherCompany%20ICM/job/icm-varnish-k8s-operator/job/master/)

### Project status: alpha
The project is in development and breaking changes can be introduced.

The purpose of the project is to provide a convenient way to deploy and manage Varnish instances in Kubernetes.

Kubernetes version `>=1.11.0` is supported.

Varnish version 6.0.2 is supported.

Full documentation can be found [here](https://pages.github.ibm.com/TheWeatherCompany/icm-varnish-k8s-operator/)

### Overview

Varnish operator manages Varnish clusters using a CustomResourceDefinition that defines a new Kind called `VarnishService`. It is conceptually similar to [Kubernetes Service](https://kubernetes.io/docs/concepts/services-networking/service/) but with a layer of Varnish instances in front of your backends. The backends are selected using a selector, just as in an ordinary Kubernetes Service.

The operator manages the whole lifecycle of the cluster: creating, deleting and keeping the cluster configuration up to date. The operator is responsible for building the VCL configuration using templates defined by the users and keeping the configuration up to date when relevant events occur (backend pod failure, scaling of the deployment, VCL configuration change).

## Features

 * [x] Basic install
 * [x] Full lifecycle support (create/update/delete)
 * [x] Automatic VCL configuration updates (using user defined templates)
 * [x] Prometheus metrics support
 * [x] Scaling
 * [x] Configurable update strategy
 * [ ] Multiple Varnish versions support
 * [ ] Autoscaling
 * [ ] Nested directors support

### Further reading

* [QuickStart](https://pages.github.ibm.com/TheWeatherCompany/icm-varnish-k8s-operator/quick-start.html)
* [Contributing](https://pages.github.ibm.com/TheWeatherCompany/icm-varnish-k8s-operator/development.html)
