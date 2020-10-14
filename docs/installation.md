# Varnish Operator 

## Installation

### Prerequisites

* Kubernetes v1.16 or newer and `kubectl` configured to communicate with your cluster
* Helm

### Configure Helm repo

```bash
$ helm repo add varnish-operator https://raw.githubusercontent.com/IBM/varnish-operator/main/helm-releases
$ helm repo update
```

### Install Varnish Operator

Create a namespace for the operator:

```bash
$ kubectl create ns varnish-operator
```

Install the operator:

```bash
$ helm install --name varnish-operator --namespace varnish-operator varnish-operator/varnish-operator
```                                                                                                                        

You should see your operator pod up and running:

```bash
$ kubectl get pods --namespace varnish-operator
NAME                              READY   STATUS              RESTARTS   AGE
varnish-operator-fd96f48f-gn6mc   1/1     Running             0          40s
```

Also, your cluster should have a new resource - `varnishcluster` (with aliases `varnishclusters` and `vc`). Now you're ready to create `VarnishCluster` resources.

To tolerate pod failures, you can run multiple replicas of the operator by specifying the `replicas` field to more than `1` in Helm value overrides. Only one of the replicas - the leader - will handle the work while others monitor if the leader is healthy. If the leader fails, one of the backup replicas will start to handle the work. That also means that if you need to check operator logs for debugging purposes you need to find the leader pod and check its logs.

See [Operator Configuration](operator-configuration.md) section for more options to configure.

## Operator Update

Since the operator is packaged into a Helm chart, the update is done by a simple `helm upgrade` command.

Note that when the operator version updates, the Varnish images get updated as well. That means that Varnish pods needs to be restarted with new configuration. As Varnish is an in-memory cache, it means losing cache data. To prevent accidental cache loss, by default, the update strategy is `OnDelete` which means the pods won't automatically get restarted. To update the pod you need to delete the pod manually, and it will come back with the new configuration. This behavior can be changed by setting the desired update strategy in the `.spec.statefulSet.updateStrategy` object. See [VarnishCluster Configuration](varnish-cluster-configuration.md) section for more details.  

## Uninstalling the Operator

Uninstallation of the chart is also done by Helm.
It deletes all created resources, including the CRD.
 
The operator uses [finalizers](https://kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/#finalizers) to perform some clean-up operations. That means you need to delete all `VarnishCluster` resources before deleting the operator. Otherwise those resources won't be deleted automatically and get stuck on your cluster.

If you happen to remove the operator first, to delete the remaining `VarnishCluster` resources you need to manually edit them and remove the finalizers:

  `kubectl patch varnishcluster <your-varnishcuster> -p '{"metadata":{"finalizers": []}}' --type=merge`

Then you need to delete the remaining objects:
 * ClusterRole named <varnishcluster-name>-varnish-clusterrole-<varnishcluster-namespace>
 * ClusterRoleBinding named <varnishcluster-name>-varnish-clusterrolebinding-<varnishcluster-namespace>