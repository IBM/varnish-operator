# Varnish As a Kubernetes Operator

I'll try to chronicle all the things I am learning as I go so that others may not have to hit as many dead ends as I have.

## Getting Started With Development

1. Have a Go environment [set up](http://sourabhbajaj.com/mac-setup/Go/README.html)
    * This primarily involves installing Go and setting a GOPATH
1. Install [`dep`](https://github.com/golang/dep)

## What is an Operator

An operator is the combination of 2 kubernetes concepts that are at the heart of how resources are managed. They are

### Resource Definitions

Literally the yaml file that describes how you want your resource to look. For instance, a yaml file specifying a deployment, and a service that should front the deployment. For the purpose of the operator, kubernetes has the concept of a CustomResourceDefinition, which allows user-defined resources to live right alongside native ones (such as service or deployment).

[here](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/) is some pretty good documentation on CustomResourceDefinitions. However, technically, the CustomResourceDefinition isn't strictly needed for some operators if they simply extend the functionality of an existing resource. For instance, if you want to manipulate Nodes, you could apply annotations onto nodes, and just watch for those annotations instead of a CustomResourceDefinition.

### Controller

When you tell Kubernetes about a resource you want (e.g. in the form of a Resource Definition), the Controller is what responds to that definition to make the resource a reality. A Controller is a constantly running piece inside Kubernetes that watches for new/updated resources, compares the desired state as described in the new Resource Definition to the actual state of the cluster, and then takes action to move the actual state closer to the desired state. As a simple example, when you scale a deployment, thus increasing/decreasing the number of desired pods, the deployment controller will react to that desire and create/destroy pods until the 2 states merge. This is formally known as a reconciliation loop. For the purpose of the operator, you will write your own controller that watches for changes to a resource, and reacts to bring the current state closer to the desired state.

Most custom controllers will take advantage of other kubernetes primitives (eg, not-custom resources like Service, Deployment, Role/Rolebinding, etc) as their primary means of control over the system.

The controller has a more nebulous definition, primarily because it can be whatever you need it to be. It's really just an application you run in the cluster that is triggered on changes to the resource it is monitoring, and reacts in some custom way. It can be written in any language where a kubernetes client sdk exists, but it is most commonly written in Go, and I will be writing the rest of this with the assumption that you are writing it in Go too.

It is highly suggested that the controller takes advantage of other kubernetes primitives (eg, not-custom resources like Service, Deployment, RBAC, etc) as much as possible, since that improves understandability. It will do that by utilizing a kubernetes client sdk, and specifically for this guide [client-go](https://github.com/kubernetes/client-go). Note that releases are versioned based on the kubernetes release, so for instance you can ask for 1.10.2, 1.9.8, etc.

### Summarize

You may want to write an operator when some or all of the below are true:

* You have logic that is hard to capture just with regular deployments/services/etc
* You need to interact directly with the kubernetes API to achieve some functionality
* You want to build a complex system that is still resilient in all the ways native kubernetes resources are
* You want to package functionality in a shareable format for others to use

Operators are a quickly evolving concept, so it's hard to find good documentation on the topic that is up-to-date. However, [the talk given by the developers of the design pattern](https://www.youtube.com/watch?v=U2Bm-Fs3b7Q) remains relevant because they discuss the concepts and motivation behind the operator.

## SDKs

While controllers are not inherently tied to any programming language, it is safe to say that Go is the predominant language used for all things related to Kubernetes, of course including Kubernetes itself. The rest of this documentation assumes that your code will be written in Go, and all SDKs discussed hereafter are for the Go language.

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

There is a "pre-alpha" tool called [operator-sdk](https://github.com/operator-framework/operator-sdk) that, in theory, makes it much easier to develop an operator. It may never leave pre-alpha, and suffice to say, at time of writing it is too half-baked to use. At the very least, be on the lookout for leaving a "pre-alpha" state.

### Kubebuilder

As of writing, [Kubebuilder](https://book.kubebuilder.io/) just hit `1.0.0`. This is the tool everyone should use. It simplifies so much of the operator writing process. It is a bit finicky to set up, but once set up it is exactly the abstraction you need to build an operator with the least amount of fuss (which is still a lot of fuss).

The remainder of this guide will discuss 2 things: How to use Kubebuider, and then at the end, some underlying implementation details from Kubebuilder that will hopefully not be relevant to users of Kubebuilder.

## Using Kubebuilder To Write An Operator

### Installation

First, follow [the installation instructions carefully](https://book.kubebuilder.io/getting_started/installation_and_setup.html). Deviating from the instructions will likely cause Kubebuilder to break at a later step, so be warned.

Do yourself a favor and don't install from master...

### The Kubebuilder Concepts

Follow [the quite good gitbook](https://book.kubebuilder.io/quick_start.html) on getting started. In fact, read through the entire gitbook. It is pretty good documentation (which is unfortunately rare in this area).

Here is a quick rundown, describing more the concepts of what's happening instead of step-by-step instructions.

#### Managers

Kubebuilder works at the top level by creating a `Manager`. The `Manager`, well, manages the lifecycle of the controller(s) that are part of the operator. You can have an arbitrary number of controllers associated with a `Manager`, although it is likely you will only have one controller for your operator.

#### (Kubebuilder's) Controller

Underneath a `Manager` are an arbitrary number of `Controller`s. `Controller`s are the same concept [as described above](#controller), meaning they watch a resource and enact a reconciliation loop on changes to that resource.

#### `kubebuilder init`

`kubebuilder init` is one of 2 commands used for project creation. You can think of this step as creating the `Manager`. It is run once per operator.

**PRO TIP**: Kubebuilder checks for the existence of a file, `hack/boilerplate.go.txt`, and it will prepend its contents to all generated Go files. If that file does not exist, it decides on a default file of

```go
/*
.
*/
```

which is incredibly annoying. You can either set the `boilerplate.go.txt` file to whatever prefix you'd like (be it authors, copyright, etc), or just create an empty file with that name to forgo any prefix.

#### `kubebuilder create api`

`kubebuilder create api` is the other of 2 commands used for project creation. You can think of this step as creating a `Controller` that will be managed by the `Manager`. It is run once for every controller your operator will have.

### What Is Generated, What Do You Write

Kubebuilder generates quite a bit for you, as well as handle much of the reconciliation loop in libraries.

Out of the box, you get scaffolding and code generation for:

* Skeleton of your custom resource definition. It uses code generation to generate the yaml file for the resource as well.
* Infrastructure for the manager/controller relationship. When you run `kubebuilder init`, the `Manager` is automatically started in the `main` function, and when you run `kubebuilder create api`, it automatically registers the `Controller` with the `Manager`.
* Convenient `Makefile` that combines many of the commands you will execute frequently, even including relevant `docker` commands. It is possibly you may need to modify the `Makefile` lightly to either customize a bit of the functionality or add new functionality (such as integration with `helm`).
* yaml generation of `Manager` StatefulSet and all RBAC for the operator. In addition, it is backed by [Kustomize](https://github.com/kubernetes-sigs/kustomize) as part of the build pipeline. See the [Kustomize](#kustomize) section for more on this tool.
* Preconfigured [code-generator](https://github.com/kubernetes/code-generator) for convenient typed clients for your CustomResourceDefinition. See the [code-generator](#code-generator) section for more on this tool, although you hopefully will never need to know more detail about it.

You will still need to fill in:

* The CustomResource definition, meaning what fields the CustomResourceDefinition should have.
* Which resources a given controller should watch for changes on (the CustomResource is already put in).
* The logic of the reconciliation loop. This is arguably the most important piece of code in the entire controller. It describes how the controller should react to changes in a watched resource.
* Good logging (by default, no logging is done on error during the reconciliation loop, which can make things hard to debug)
* External state to the operator. Meaning, if any of the state needs to be injected at installation time, that will need to be configured.

## Background Information

Kubebuilder handles much of the work under the hood, but sometimes it might be worth knowing what's going on under there. What follows is an incomplete list of concepts that Kubebuilder is covering

### Package Managers

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

After following those steps, any incoming CRD for which you have defined defaults should inherit those defaults if a value isn't explicitly set for them. If you need to add any more defaults, or modify existing ones, just edit the `defaults.go` file, and run

```sh
go generate ./pkg/...
```

again.

### Logging

Logging implementation uses [zap](https://github.com/uber-go/zap) logging framework.

Logging configuration accepts standard log levels in string format - `error`, `warn`, `info` and `debug`. Default log level is set to `info`, but can be configured by setting `LOG_LEVEL` environment variable to desired log level.

To make logs produced by operator easy to parse and store in external logging systems, `json` encoder is used by default. Since it is not very human-friendly, it is possible to switch logging encoder to `console` encoder by setting `LOG_FORMAT` environment variable value to `console`. 

When Varnish operator is deployed using Helm, logging configuration is handled by setting values `operator.loglevel` and `operator.logformat` in [values.yaml](./varnish-operator/values.yaml)  accordingly.

During local development, `make run` command sets logging encoder to `console`. To increase logging verbosity, one can do so by setting desired log level using environment variable: 
```bash
LOG_LEVEL=debug make run
```

## Deploying Your Kubebuilder project

Out of the box, Kubebuilder has an integration with [Kustomize](https://github.com/kubernetes-sigs/kustomize). It bills itself as letting you "customize raw, template-free YAML files for multiple purposes, leaving the original YAML untouched and usable as is." You can learn more about Kustomize from the link. It is a very simple tool, and all of the documentation takes just a few minutes to read through.

Kustomize fills a similar space to [Helm](https://helm.sh), which is very much a templating engine, as well as various other things, like a deployment platform.

The ICM team has standardized on Helm, so you will need to make Kustomize and Helm work together.

### Kubebuilder + Helm

When trying to get these 2 tools to work together, you will quickly find that Kustomize actually parses the yaml files in order to generate its final result, while Helm introduces non-yaml syntax in the form of go templates. Thus, if you try to use Kustomize on a yaml file that has Helm template elements in it, Kustomize will blow up. This is the main limitation you will need to work around to get these 2 tools to work together. It is actually quite straightforward to work in the opposite direction -- meaning you can generate valid yaml files by running Helm, and then pipe those files to Kustomize, which can do its own replacements.

Since Kubebuilder generates yaml files intended to then be passed to Kustomize, you're only left with the option to move from Kustomize into Helm, meaning the first pass cannot have any Helm templates in them. However, that only applies to file that the Kubebuilder project generates every time you call `make`, which are:

* One CRD for every controller you have
* RBAC files, namely a Role and Rolebinding

Which means there are some files which are not generated every time, and they are:

* All files related to the manager, meaning the namespace the manager resides in, its statefulset, and its deployment. All 3 of these are in the same file under `config/manager`.

However, note that there are 2 things that still affect the manager files: a Kustomize patch file, `manager_image_patch.yaml`, which updates the statefulset with the correct version of the manager; and the namespace+name prefix fields in the `kustomization.yaml` file.

One possible workflow, given these facts, is to first generate the CRD and RBAC yaml files through the build process, run Kustomize on them, and move the output unedited directly into a helm charts directory. Also, permanently move the Manager file to that directory. The Manager file can have any helm templating needed, which should include the correct version, a name prefix, and namespace which should all be parsed from the `kustomization.yaml` file directly at build time, to make up for the loss of Kustomize. In this way, the helm chart will have templating as part of the manager, which crucially allows input into how the controllers act at installation time, as opposed to just at compile time.

**NOTE**: Make sure all resources are broken up into individual files (i.e. not all in the same file separated by `--`), since [Helm has an issue otherwise](https://github.com/helm/helm/issues/3785)).
