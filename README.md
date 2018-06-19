# Varnish As a Kubernetes Operator

I'll try to chronicle all the things I am learning as I go so that others may not have to hit as many dead ends as I have.

## Getting Started With Development

1. Have a Go environment [set up](http://sourabhbajaj.com/mac-setup/Go/README.html)

    * This primarily involves installing Go and setting a GOPATH

1. Install [`dep`](https://github.com/golang/dep)

## operator-sdk

There is a "pre-alpha" tool called [operator-sdk](https://github.com/operator-framework/operator-sdk) that, in theory, makes it much easier to develop an operator. It may never leave pre-alpha, and suffice to say, at time of writing it is too half-baked to use. At the very least, be on the lookout for leaving a "pre-alpha" state.

Without the operator-sdk, you are left with rolling your own operator. The documentation to do that is very thin, but in truth operators are more of a design pattern, so understanding the pattern is really the important part.

## What is an Operator

An operator is a concept in kubernetes that is at the heart of how resources are managed. For instance, when you upload a yaml file specifying a Deployment resource, you are actually telling the built-in Deployment operator to figure out how to actually stand up the pods/containers.

To take a step back for a second, an operator is primarily made of 2 pieces:

**Resource Definition**: This is literally the yaml file that you give kubernetes describing how you want your resource to look. For instance, a service.yaml, deployment.yaml, etc. Kubernetes has the ability to accept a CustomResourceDefinition, which allows a user-defined resource with an arbitrary spec. When writing your own operator, you'll likely need to specify a CustomResourceDefinition.

**Controller**: This is the running piece inside kubernetes that watches for new/updated resources, compares the desired state as specified in the resource and the actual state of the cluster, and to move the actual state towards the desired state. This is known as a reconciliation loop.

Technically, the Resource Definition isn't strictly needed for some operators, if they simply extend the functionality of an existing resource. For instance, If you want to manipulate some Nodes, you could just use annotations that a controller is watching for, and take action only on that.

### Custom Resource Definitions

The Custom Resource Definition, or CSR, allows you to make an arbitrary resource that `kubectl` can see and act on. [here's](https://kubernetes.io/docs/tasks/access-kubernetes-api/extend-api-custom-resource-definitions/) some pretty good documentation on how they work.

### Controller

The controller has a more nebulous definition, primarily because it can be whatever you need it to be. It's really just an application you run in the cluster that is triggered on changes to the resource it is monitoring, and reacts in some custom way. It can be written in any language where a kubernetes client sdk exists, but it is most commonly written in Go, and I will be writing the rest of this with the assumption that you are writing it in Go too.

It is highly suggested that the controller takes advantage of other kubernetes primitives (eg, not-custom resources like Service, Deployment, RBAC, etc) as much as possible, since that improves understandability. It will do that by utilizing a kubernetes client sdk, and specifically for this guide [client-go](https://github.com/kubernetes/client-go). Note that releases are versioned based on the kubernetes release, so for instance you can ask for 1.10.2, 1.9.8, etc.

### Summarize

To build an operator, you need some resource that you need to monitor, be it a Custom Resource or a built-in resource, and a controller that does the monitoring, and enacts a reconciliation loop to move the state of the cluster to look more like the resource.

If you want to understand operators better, there is [a talk given by the developers of the design pattern](https://www.youtube.com/watch?v=U2Bm-Fs3b7Q) that was very enlightening. I'm sorry that this is the best documentation available...

## Useful Tools

This section will be more of a grab-bag of useful tools I found along the way. You will probably need most of them, but I will point out when I think a tool is optional, and why.

### Package Managers

Unfortunately, Go has the distiction of having about 10 different package managers. You may find trying to navigate that really confusing. Fortunately, the community has finally convened on one package manager. But unforunately, that does not mean that all libraries play nicely with this manager. So, sometimes you will still need to make a decision on which package manager to pick. But fortunately, it's only a decision between 2 different managers. Those managers are:

[**dep**](https://github.com/golang/dep): The "official" package manager. However, it is new, so it does not have recursive dependency management. Sometimes, it is still possible to include all dependencies manually, but that can be confusing at times, in which case, you would consider using...

[**glide**](https://github.com/Masterminds/glide): Before `dep`, `glide` was the biggest package manager, and had gotten as far as supporting recursive dependencies.

As it turns out, the `client-go` library has recursive dependencies, so one might think that `glide` is the right choice. However, all examples of operators I have seen use `dep` anyway, and manually specify the needed dependencies. I have tried both and they both work, so take it for what it's worth. For the rest of this guide, I will try to use `dep`.

### Helm \*\*Possibly Optional\*\*

[Helm](https://helm.sh) lets you use templates with your deployment, which is almost always very nice to have because kubernetes yaml files are often repetitive, and this cuts down on errors. It comes with its own complications and idiosyncracies, however, so for a first pass you might try without using Helm to reduce overall complexity, then add it in later.

This is talked about elsewhere. I will link it when I find that.

### client-go

[This library](https://github.com/kubernetes/client-go) is the heart of how the operator works. It allows the operator to communicate with kubernetes, which allows you to monitor the resource and enact changes on the cluster in response to changes in the resource. It is...poorly documented, but suffice to say it roughly mirrors the kubernetes [HTTP API](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.10), with a few added features. There are some very basic [examples](https://github.com/kubernetes/client-go/tree/master/examples) that show how you can set up the client, but it does not go very much beyond that.

#### Informer

There is an incredibly important concept in the client go that is poorly referenced called an _Informer_ that basically acts as a caching layer to the kubernetes API. It is explained in [this youtube video](https://www.youtube.com/watch?v=U2Bm-Fs3b7Q) and in [this random page of documentation](https://github.com/kubernetes/sample-controller/blob/master/docs/controller-client-go.md).

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

### Code Generation

There is a decent amount of code that is very rote in nature, but is required for your CRD definitions. Primarily, these definitions must be an instance of `runtime.Object`, which requires a `DeepCopy` method to exist. While you could write this definition yourself, there is a [code generation library](https://github.com/kubernetes/code-generator) that can handle it for you. It can additionally handle integrations with lister, cacher, and client.

Follow along with [this guide](https://blog.openshift.com/kubernetes-deep-dive-code-generation-customresources/) to use code generation.

There are a few finicky parts to the code generation tool that are worth noting:

#### Gopkg.toml Dependency

When including the code-generation library as a dependency with `dep`, it becomes quickly evident that there is no actual code dependency on the library, so `dep` will not download anything. In fact, `dep` will complain that it is not being used when running `dep ensure`. At time of writing, there is no "supported" way to fix this, but there is a hacky and obtuse workaround possible. [Following this github issue](https://github.com/kubernetes/sample-controller/issues/6), it is possible to add special `prune` settings that force `dep` to download the files:

```toml
[[constraint]]
  name = "k8s.io/code-generator"
  version = "kubernetes-1.10.2"

[prune]
  non-go = true
  go-tests = true
  unused-packages = true

  [[prune.project]]
    name = "k8s.io/code-generator"
    unused-packages = false
    non-go = false
    go-tests = false

  [[prune.project]]
    name = "k8s.io/gengo"
    unused-packages = false
    non-go = false
    go-tests = false
```

* The first block `[prune]` is the default settings for all packages
* The second block and third block `[[prune.project]]` configure these two packages to download all packages, whether they are used or not
* It is unclear why the third block is needed. The person in the github explains it, but does so insufficiently... suffice to say, without the third block, code generation will not work

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

* The first two lines should always be the same, since this script SHOULD be sitting inside a top-level folder called "hack". For a full project layout, [see the appropriate section]()
* `<go-path-to-project>` is the full path to the project folder starting at $GOPATH. For instance, if your project is at `$GOPATH/icm.ibm.com/myproject`, then use `icm.ibm.com/myproject`
* `<folder-under-apis>` is the folder name within `pkg/apis`. For instance, if it is `pkg/apis/icm.ibm.com/`, then use `icm.ibm.com`
* `<version>` is the version used, and should also match the folder that is inside `<folder-name-under-apis>`. For instance, if it is `pkg/apis/icm.ibm.com/v1beta1/`, use `v1beta1`

#### verify-codegen.sh

This script can be copied exactly as it exists in the sample project.