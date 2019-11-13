# Varnish As a Kubernetes Operator

I'll try to chronicle all the things I am learning as I go so that others may not have to hit as many dead ends as I have.

## Kubebuilder Did It Better

Kubebuilder [covers similar topics](https://book.kubebuilder.io) to much of this documentation. It is highly recommended that you read through all of Kubebuilder's docs, and then return here and only visit sections missed, or are specific to this repo 

## Getting Started With Development

Follow along with [Kubebuilder's Quick Start Guide](https://book.kubebuilder.io/quick_start.html)

## What is an Operator

An operator is the combination of 2 kubernetes concepts that are at the heart of how resources are managed. They are

### Resource Definitions

Literally the yaml file that describes how you want your resource to look. For instance, a yaml file specifying a deployment, and a service that should front the deployment. For the purpose of the operator, kubernetes has the concept of a CustomResourceDefinition, which allows user-defined resources to live right alongside native ones.

[here](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/) is some pretty good documentation on CustomResourceDefinitions. However, technically, the CustomResourceDefinition isn't strictly needed for some operators if they simply extend the functionality of an existing resource. For instance, if you want to manipulate Nodes, you could apply annotations onto nodes, and just watch for those annotations instead of a CustomResourceDefinition.

Also see [Kubebuilder docs around Resources](https://book.kubebuilder.io/basics/what_is_a_resource.html)

### Controller

When you tell Kubernetes about a resource you want (e.g. in the form of a Resource Definition), the Controller is what responds to that definition to make the resource a reality. A Controller is a constantly running piece inside Kubernetes that watches for new/updated resources, compares the desired state as described in the new Resource Definition to the actual state of the cluster, and then takes action to move the actual state closer to the desired state. As a simple example, when you scale a deployment, thus increasing/decreasing the number of desired pods, the deployment controller will react to that desire and create/destroy pods until the 2 states merge. This is formally known as a reconciliation loop. For the purpose of the operator, you will write your own controller that watches for changes to a resource, and reacts to bring the current state closer to the desired state.

Most custom controllers will take advantage of other kubernetes primitives (eg, not-custom resources like Service, Deployment, Role/Rolebinding, etc) as their primary means of control over the system.

The controller has a more nebulous definition, primarily because it can be whatever you need it to be. It's really just an application you run in the cluster that is triggered on changes to the resource it is monitoring, and reacts in some custom way. It can be written in any language since kubernetes has an HTTP API, but it is most commonly written in Go, and I will be writing the rest of this with the assumption that you are writing it in Go too.

See [Kubebuilder's take on controllers as well](https://book.kubebuilder.io/basics/what_is_a_controller.html)

### Summarize

You may want to write an operator when some or all of the below are true:

* You have logic that is hard to capture just with regular deployments/services/etc
* You need to interact directly with the kubernetes API to achieve some functionality
* You want to build a complex system that is still resilient in all the ways native kubernetes resources are
* You want to package functionality in a shareable format for others to use

Operators are a quickly evolving concept, so it's hard to find good documentation on the topic that is up-to-date. However, [the talk given by the developers of the design pattern](https://www.youtube.com/watch?v=U2Bm-Fs3b7Q) remains relevant because they discuss the concepts and motivation behind the operator.

Also [Kubebuilder has built some very good docs around its approach to operators](https://book.kubebuilder.io)

## Kubebuilder

This project uses [Kubebuilder](https://book.kubebuilder.io/), and it is highly recommended that you read through all of their documentation.

This guide will NOT cover how to use Kubebuilder, since its documentation is very good. It will, however, discuss other frameworks, and lower-level details.

## Other SDKs

### client-go

[client-go](https://github.com/kubernetes/client-go) is the official Go client for Kubernetes, acting primarily as a thin wrapper around the API. It inherits some of the core concepts that are a part of core Kubernetes as well. It is a good starting point when exploring the Kubernetes API, and is the go-to tool when writing any non-operator projects that must interact with Kubernetes programmatically.

The most important concept it exposes is that of Informers:

#### Informers

The _Informer_ basically acts as a caching layer to the kubernetes API. It is explained in [this random page of documentation](https://github.com/kubernetes/sample-controller/blob/master/docs/controller-client-go.md) and in [this youtube video](https://www.youtube.com/watch?v=U2Bm-Fs3b7Q).

For a practical example of what that is, see this StackOverflow question where some [random dude writes some code](https://stackoverflow.com/questions/40975307/how-to-watch-events-on-a-kubernetes-service-using-its-go-client).

Reproduced here:

Q: How to watch events on a kubernetes service using its go client

A: this can be done like this:

```go
package main

import (
    "fmt"
    "flag"
    "time"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/pkg/api/v1"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/tools/cache"
    "k8s.io/client-go/pkg/fields"
)

var (
    kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
)

func main() {
    flag.Parse()
    config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
        panic(err.Error())
    }
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }

    watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), "services", v1.NamespaceDefault,
        fields.Everything())
    store, controller := cache.NewInformer(
        watchlist,
        &v1.Service{},
        time.Second * 0,
        cache.ResourceEventHandlerFuncs{
            AddFunc: func(obj interface{}) {
                fmt.Printf("service added: %s \n", obj)
            },
            DeleteFunc: func(obj interface{}) {
                fmt.Printf("service deleted: %s \n", obj)
            },
            UpdateFunc:func(oldObj, newObj interface{}) {
                fmt.Printf("service changed \n")
            },
        },
    )
    stop := make(chan struct{})
    go controller.Run(stop)
    for{
        time.Sleep(time.Second)
    }
}
```

`NewListWatchFromClient` wraps around the _Informer_ for a cleaner usage, but it essentially caches all requests made against the Kubernetes API. In the example above, `store` has functions like `List` and `Get` that will either respond from cache or actually execute the `List`/`Get`. `controller` executes a `Watch` for changes in the resource, as well as a reload of the cache every `resync` period (usually 1 minute).

Unfortunately it is still possibly worthwhile to just dive into the source code and poke around.

Despite that, I would not watch resources without this feature.

### operator-sdk

**NOTE**: After taking a look again at operator-sdk, it appears to have been completely overhauled to use the same underlying Go library as Kubebuilder. They also seem to have added real documentation. This is worth looking into.
 
## Lower-Level Information

Kubebuilder handles much of the work under the hood, but sometimes it might be worth knowing what's going on under there. What follows is an incomplete list of concepts that Kubebuilder is covering

### Package Managers

**NOTE**: this section is no longer true. Seems like the community is now shifting to `vgo` aka the package manager integrated with the `go` binary. As of writing, Kubebuilder still uses `dep`, but may one day switch to `go`.

Unfortunately, Go has the distiction of having about 10 different package managers. You may find trying to navigate that really confusing. Fortunately, the community has finally convened on one package manager. But unforunately, that does not mean that all libraries play nicely with this manager. But fortunately, it's only a decision between 2 different managers. Those managers are:

[**dep**](https://github.com/golang/dep): The "official" package manager. However, it is new, so it does not have recursive dependency management. Sometimes, it is still possible to include all dependencies manually, but that can be confusing at times. Kubebuilder uses `dep` regardless, and generates all the dependencies it needs.

[**glide**](https://github.com/Masterminds/glide): Before `dep`, `glide` was the biggest package manager, and had gotten as far as supporting recursive dependencies.

#### Gopkg.toml Dependency

When including the code-generation library as a dependency with `dep`, it becomes quickly evident that there is no actual code dependency on the library, so `dep` will not download anything. In fact, `dep` will complain that it is not being used when running `dep ensure`. To fix this, you can "require" that certain packages be included in the dep build. the `required` field must go before any `[[constraint]]` or `[[override]]` field, so it is usually safest to just make it the first line in the file:

```toml
required = [
  "k8s.io/code-generator/cmd/client-gen",
  "k8s.io/code-generator/cmd/conversion-gen",
  "k8s.io/code-generator/cmd/deepcopy-gen",
  "k8s.io/code-generator/cmd/defaulter-gen",
  "k8s.io/code-generator/cmd/informer-gen",
  "k8s.io/code-generator/cmd/lister-gen",
]

[[constraint]]
  name = "k8s.io/code-generator"
  version = "kubernetes-1.10.2"
```

* The first block `[prune]` is the default settings for all packages
* The second block and third block `[[prune.project]]` configure these two packages to download all packages, whether they are used or not
* It is unclear why the third block is needed. The person in the github explains it, but does so insufficiently... suffice to say, without the third block, code generation will not work

### Code Generation

There is a decent amount of code that is very rote in nature, but is required for your CRD definitions. Primarily, these definitions must be an instance of `runtime.Object`, which requires a `DeepCopy` method to exist. While you could write this definition yourself, there is a [code generation library](https://github.com/kubernetes/code-generator) that can handle it for you. It can additionally handle integrations with lister, cacher, and client.

Follow along with [this guide](https://blog.openshift.com/kubernetes-deep-dive-code-generation-customresources/) to use code generation.

There are a few finicky parts to the code generation tool that are worth noting:

#### update-codegen.sh

This file really wants a particular format. It should look exactly like

```sh
#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(ls -d -1 $SCRIPT_ROOT/vendor/k8s.io/code-generator 2>/dev/null)}

$CODEGEN_PKG/generate-groups.sh <generators> \
<go-path-to-project>/pkg/client <go-path-to-project>/pkg/apis \
<folder-under-apis>:<version>
```

* The first two lines should always be the same, since this script SHOULD be sitting inside a top-level folder called "hack".
* `<go-path-to-project>` is the full path to the project folder starting at $GOPATH. For instance, if your project is at `$GOPATH/icm.ibm.com/myproject`, then use `icm.ibm.com/myproject`
* `<folder-under-apis>` is the folder name within `pkg/apis`. For instance, if it is `pkg/apis/icm.ibm.com/`, then use `icm.ibm.com`
* `<version>` is the version used, and should also match the folder that is inside `<folder-name-under-apis>`. For instance, if it is `pkg/apis/icm.ibm.com/v1beta1/`, use `v1beta1`

#### verify-codegen.sh

This script basically just calls `update-codegen.sh` and compares the output with the initial code.

### Defaults

Kubernetes has a built-in mechanism of setting defaults for resources, which can be very useful when creating your CRD. However, documentation for it is non-existent, so here are the steps required to have defaults for incoming CRDs:

1. In the same folder as your `<CRD-name>_types.go` (or `types.go` if not using Kubebuilder), create a `defaults.go` file ([here is the file in this project](pkg/apis/icm/v1alpha1/defaults.go)). In this file, create functions that will define defaults by type for your CRD. They MUST have the name template of `SetDefaults_<type>(in *<type>)`. For example:
    
    ```go
    SetDefaults_TypeWithDefaults(in *TypeWithDefaults) {
        if (in.FieldThatNeedsDefault == "") {  
            in.FieldThatNeedsDefault = "default"
        }
    }
    ```
1. open the `register.go` file at the same level as the `defaults.go` file, and make sure a comment

    ```go
    // +k8s:defaulter-gen=TypeMeta
    ```

   exists. This informs the code generation tool (discussed below) to create an overall default function for all types that inherit the `TypeMeta` struct. That will likely be just the CRD definition and the `<CRD-type>List` type located in the `types.go` file (for example, `MyCRD` and `MyCRDList`).
1. In the `defaults.go` file, add a comment at the top:

    ```go
    //go:generate go run ../../../../vendor/k8s.io/code-generator/cmd/defaulter-gen/main.go -O zz_generated.defaults -i . -h ../../../../hack/boilerplate.go.txt
    ```

   The number of `../` may differ, depending on what relative path is necessary to reach the `vendor` and `hack` folder. This informs go to run the [`defaulter-gen`](https://godoc.org/k8s.io/gengo/examples/defaulter-gen) code generation script, which will look for the `+k8s:defaulter-gen` comment and any `SetDefaults_*` templated functions to create `SetObjectDefaults_*` for all types that inherit `TypeMeta` into a `zz_generated.defaults.go` file.
1. Run the code generation by calling

    ```sh
    go generate ./pkg/...
    ```

   which will look through all source files in `pkg` for a `//go:generate` comment and execute the command there.
1. In the (now generated) `zz_generated.defaults.go` file ([here is the file in this repo](pkg/apis/icm/v1alpha1/zz_generated.defaults.go)), there will be a `RegisterDefaults` function, which is the hook used to tell Kubernetes how to use the generated default functions. If using Kubebuilder, all you need to do is open the `addtoscheme_<group-name>_<version>.go` file inside `pkg/apis` and add the `RegisterDefaults` function to the `AddToSchemes` slice. For example:

    ```go
    AddToSchemes = append(AddToSchemes, v1alpha1.SchemeBuilder.AddToScheme, v1alpha1.RegisterDefaults)
    ```

   If you are not using Kubebuilder, you'll need to inform the `scheme.Scheme` for your controller of the defaults by running

    ```go
    RegisterDefaults(<your-scheme-here>)
    ```

After following those steps, your scheme will be able to correctly handle defaulting defined types using the `Default` function. For instance, at the top of your reconcile function, after `Get`ting the current CRD `instance`, you may fill in defaults for that `instance` by running
 
 ```go
 scheme.Default(instance)
```

and your CRD should now inherit those defaults if a value isn't explicitly set for them.

If you need to add any more defaults, or modify existing ones, just edit the `defaults.go` file, and run

```sh
go generate ./pkg/...
```

again.

### Logging

**NOTE**: Kubebuilder is now providing a wrapper around `zap` that should probably be investigated and used, in lieu of `zap`.

Logging implementation uses [zap](https://github.com/uber-go/zap) logging framework.

Logging configuration accepts standard log levels in string format - `error`, `warn`, `info` and `debug`. Default log level is set to `info`, but can be configured by setting `LOG_LEVEL` environment variable to desired log level.

To make logs produced by operator easy to parse and store in external logging systems, `json` encoder is used by default. Since it is not very human-friendly, it is possible to switch logging encoder to `console` encoder by setting `LOG_FORMAT` environment variable value to `console`. 

When Varnish operator is deployed using Helm, logging configuration is handled by setting values `operator.loglevel` and `operator.logformat` in [values.yaml](./varnish-operator/values.yaml)  accordingly.

During local development, `make run` command sets logging encoder to `console`. To increase logging verbosity, one can do so by setting desired log level using environment variable:

```bash
LOG_LEVEL=debug make run
```

## Publishing Kubernetes events

Kubebuilder has a built in mechanism for publishing Kubernetes events. To use it you need to create an event recorder with help of the manager. It can be done on controller creation so you can save and use it later in your reconciliation logic:

 ```go
package myservice

import (
	...
	"k8s.io/client-go/tools/record"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	...
)
// ReconcileMyService reconciles MyService
type ReconcileMyService struct {
	client.Client
	scheme *runtime.Scheme
	events *record.EventRecorder
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileMyService{
		Client: mgr.GetClient(),
		scheme: mgr.GetScheme(),
		events: mgr.GetRecorder("my-service"),
	}
}
 ``` 

 then you can use it your Reconcile() function:
 
 ```go
 func (r *ReconcileMyService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
 	instance := &icmv1alpha1.MyService{}
 	r.events.Event(instance, "Normal", "ReconciliationEvent", "MyService reconciliation is started")
 	err := r.Get(context.TODO(), request.NamespacedName, instance)
 		...
 }
 ```

* **component** will appear in the "From" field of the events and is the component that creates this events
* **type** must be **Normal** or **Warning**
* **reason** is the reason why the event is generated. It should be in UpperCamelCase format and be short. It will be used by automation as an identifier.
* **message** the description of the event that happened

Here's an example of the events that you will see in your `kubectl describe <your-component>` output:

```
 ...
    Selector:
      App:  HttPerf
Events:
  Type    Reason           Age                   From             Message
  ----    ------           ----                  ----             -------
  Warning  DeploymentError  5m7s (x32 over 4h2m)  component-name  Could not update deployment. Error: some error
  
```

## Deploying Your Kubebuilder project

**NOTE**: Kustomize is now part of the `kubectl` command line tool, meaning it is apparently becoming a standard part of kubernetes deployments. We should investigate how much and where to use Kustomize generally

Out of the box, Kubebuilder has an integration with [Kustomize](https://github.com/kubernetes-sigs/kustomize). It bills itself as letting you "customize raw, template-free YAML files for multiple purposes, leaving the original YAML untouched and usable as is." You can learn more about Kustomize from the link. It is a very simple tool, and all of the documentation takes just a few minutes to read through.

Kustomize fills a similar space to [Helm's templating engine](https://helm.sh).

The ICM team has standardized on Helm, so you will need to make Kustomize and Helm work together.

### Kustomize + Helm

When trying to get these 2 tools to work together, you will quickly find that Kustomize actually parses the yaml files in order to generate its final result, while Helm introduces non-yaml syntax in the form of go templates. Thus, if you try to use Kustomize on a yaml file that has Helm template elements in it, Kustomize will blow up. This is the main limitation you will need to work around to get these 2 tools to work together. It is actually quite straightforward to work in the opposite direction -- meaning you can generate valid yaml files by running Helm, and then pipe those files to Kustomize, which can do its own replacements.

Since Kubebuilder generates yaml files intended to then be passed to Kustomize, you're only left with the option to move from Kustomize into Helm, meaning the first pass cannot have any Helm templates in them. However, that only applies to file that the Kubebuilder project generates every time you call `make`, which are:

* One CRD for every controller you have
* RBAC files, namely a Role and Rolebinding

Which means there are some files which are not generated every time, and they are:

* All files related to the manager, meaning the namespace the manager resides in, its statefulset, and its deployment. All 3 of these are in the same file under `config/manager`.

However, note that there are 2 things that still affect the manager files: a Kustomize patch file, `manager_image_patch.yaml`, which updates the statefulset with the correct version of the manager; and the namespace+name prefix fields in the `kustomization.yaml` file.

One possible workflow, given these facts, is to first generate the CRD and RBAC yaml files through the build process, run Kustomize on them, and move the output unedited directly into a helm charts directory. Also, permanently move the Manager file to that directory. The Manager file can have any helm templating needed, which should include the correct version, a name prefix, and namespace which should all be parsed from the `kustomization.yaml` file directly at build time, to make up for the loss of Kustomize. In this way, the helm chart will have templating as part of the manager, which crucially allows input into how the controllers act at installation time, as opposed to just at compile time.

**NOTE**: Make sure all resources are broken up into individual files (i.e. not all in the same file separated by `--`), since [Helm has an issue otherwise](https://github.com/helm/helm/issues/3785)).

### GenerateName usage

Kubernetes requires that only one object of a given kind can have a given name. If it's hard to ensure name uniqueness, it is possible to ask Kubernetes to generate unique names.

To do so, instead of specifying the name of the object explicitly, you set the `generateName` field for your resource:

```go
serviceAccount := &v1.ServiceAccount{
    ObjectMeta: metav1.ObjectMeta{
        GenerateName: "service-account-prefix-",
        Namespace:    instance.Namespace,
    }
}
```

Kubernetes will use that value as a prefix for the new name. Suffix will be a unique string. 

There is a length limit on names, including prefix for generated names, so it may be possible that you will need to truncate the prefix value before assigning it to `GenerateName`. The docs describe the rules for limits as follows:

>By convention, the names of Kubernetes resources should be up to maximum length of 253 characters and consist of lower case alphanumeric characters, -, and ., but certain resources have more specific restrictions.

Some limits are described in more details for specific resources (e.g. Labels and Selectors) but most of them are not so be prepared to find them out only after hitting them.

Also, if the concatenation of `GenerateName` value and unique suffix exceeds the limitation, Kubernetes will truncate the prefix to fit in the limit. 

Downside of using `GenerateName` is that in the code you need to save the generated name somewhere if you later need to refer to that object. 

### Keeping Varnish Stable

Kubernetes is built on the premise that its runnable environments are ephemeral, meaning they can be created or deleted at will, with little to no effect on the overall system. In the case of Varnish, which is purely an in-memory caching layer, deleting and creating instances all the time would cause the cache to perform very poorly. Thus, there is a need to keep Varnish stable, ie tell Kubernetes that these particular runnable environments should _not_ be treated as ephemeral.

Kubernetes does not provide this functionality out of the box, but you can trick it into approximating this behavior, and that is through the concepts of guaranteed resources and affinities.

### Guaranteed Resources

The way that Kubernetes manages deployed pods on nodes is through monitoring the resources that a pod is using. Specifically, it uses the `limits` and `requests` values for `cpu` and `memory` to determine how much resources to give a pod, and when it might be OK to reschedule a pod somewhere else (namely, if a node is running out of resources and some pods are using more resources than requested). For a detailed breakdown of what `limits` and `requests` mean, [see the Kubernetes documentation on QoS](https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/). In QoS parlance, you want the Varnish nodes to be a "Guaranteed" QoS. In short, you want to always set the `limits` and `requests` fields, and you want `limits` and `requests` to be identical.

### Affinities

[Kubernetes allows control](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#inter-pod-affinity-and-anti-affinity-beta-feature) on where pods get deployed based on labels associated with pods and nodes. For instance, you can configure pods of the same deployment to repel each other, meaning new pods entering the deployment will try to avoid nodes that already have a pod of that type. That way, you if any one node goes down, it will only take a single pod with it. Likewise, you can configure pods to be attracted to each other, for colocation that could decrease latency between pods. Note that reading through the above linked documentation is valuable, as it goes into limitations to affinities, as well as deeply explains how they work and when to use them.

For the purposes of this Varnish deployment, you will most likely want to configure a [pod anti-affinity](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#never-co-located-in-the-same-node) so that each pod of the varnish deployment is on a different node. An example of what that might look like is in the [example annotated yaml file](/config/samples/icm_v1alpha1_varnishcluster.yaml) under `spec.deployment.affinity`.

