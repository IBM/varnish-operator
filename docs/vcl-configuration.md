# VCL files configuration

TODO not ready doc, just moved the relevant section from README.md here.

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
