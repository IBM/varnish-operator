# Architecture

NOT READY SECTION. JUST MOVED THE OLD RELEVANT DOCS HERE


### Components

**CustomResourceDefinition**: the actual "VarnishService" resource, that acts in the same way that a Service resource does, except with an added Varnish layer between the Service and the Deployment it backs. You would define a resource of Kind "VarnishService", and specify all the regular specs for a Service, plus some new fields that control how many Varnish instances you want, how much memory/cpu they get, and other relevant information for the Varnish cluster.

**Controller**: The controller is an application deployed into your cluster that knows how to react to the VarnishService CustomResource. Meaning, this application watches for new or changed VarnishServices and handles the actual underlying infrastructure. That means it must be running at all times in the cluster, although it lives in its own namespace away from your application.

**Varnish**: Of course, there will be a cluster of Varnishes. Currently, only Varnish version 6.1.1 is supported, but more versions will be supported in the future.

**K-Watcher**: K-watcher is a sidecar to every Varnish node that monitors the backends and updates the Varnish VCL whenever it notices a change to the pod IP addresses.


## Keeping Varnish Stable

Kubernetes is built on the premise that its runnable environments are ephemeral, meaning they can be created or deleted at will, with little to no effect on the overall system. In the case of Varnish, which is purely an in-memory caching layer, deleting and creating instances all the time would cause the cache to perform very poorly. Thus, there is a need to keep Varnish stable, ie tell Kubernetes that these particular runnable environments should _not_ be treated as ephemeral.

Kubernetes does not provide this functionality out of the box, but you can trick it into approximating this behavior, and that is through the concepts of guaranteed resources and affinities.

### Guaranteed Resources

The way that Kubernetes manages deployed pods on nodes is through monitoring the resources that a pod is using. Specifically, it uses the `limits` and `requests` values for `cpu` and `memory` to determine how much resources to give a pod, and when it might be OK to reschedule a pod somewhere else (namely, if a node is running out of resources and some pods are using more resources than requested). For a detailed breakdown of what `limits` and `requests` mean, [see the Kubernetes documentation on QoS](https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/). In QoS parlance, you want the Varnish nodes to be a "Guaranteed" QoS. In short, you want to always set the `limits` and `requests` fields, and you want `limits` and `requests` to be identical.

### Affinities

[Kubernetes allows control](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#inter-pod-affinity-and-anti-affinity-beta-feature) on where pods get deployed based on labels associated with pods and nodes. For instance, you can configure pods of the same deployment to repel each other, meaning new pods entering the deployment will try to avoid nodes that already have a pod of that type. That way, you if any one node goes down, it will only take a single pod with it. Likewise, you can configure pods to be attracted to each other, for colocation that could decrease latency between pods. Note that reading through the above linked documentation is valuable, as it goes into limitations to affinities, as well as deeply explains how they work and when to use them.

For the purposes of this Varnish deployment, you will most likely want to configure a [pod anti-affinity](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#never-co-located-in-the-same-node) so that each pod of the varnish deployment is on a different node. An example of what that might look like is in the [example annotated yaml file](/config/samples/icm_v1alpha1_varnishservice.yaml) under `spec.deployment.affinity`.

