# VCL files configuration

Varnish VCL files are stored in a ConfigMap. Each entry in the ConfigMap corresponds to a file with VCL configuration.

There are 2 fields relevant to configuring the `VarnishService` for VCL code, in `.spec.vclConfigMap` object:

* **name**: This is a required field and tells the VarnishService the name of the ConfigMap that contains/will contain the VCL files
* **entrypointFile**: The name of the file that acts as the entrypoint for Varnish. This is the name of the file that will be passed to the Varnish executable. If `entrypointFile` is templated (ends in `.tmpl`), exclude the `.tmpl` extension. For example, if ConfigMap has file `mytemplatedfile.vcl.tmpl`, set `entrypointFile: mytemplatedfile.vcl` as the generated file will omit the extension.

If a ConfigMap does not exist on VarnishService creation, the operator will create one and populate it with a default `backends.vcl.tmpl` and `default.vcl`. Their behavior is as follows:

* `backends.vcl.tmpl`: collect all backends into a single round-robin director
* `default.vcl`:
  * respond to `GET /heartbeat` checks with a 200
  * respond to `GET /liveness` checks with a 200 or 503, depending on healthy backends
  * respond to all other requests normally, caching all non-404 responses
  * hash request based on url
  * add `X-Varnish-Cache` header to response with "HIT" or "MISS" value, based on presence in cache
    
### Writing a Templated VCL File

The template file is a regular VCL file, with the addition of [Go templates](https://golang.org/pkg/text/template). This is because there is no way to know the backend's IP addresses at startup, so they must be injected at runtime. Also they can change over time if the backends get rescheduled by Kubernetes. 

The file is considered templated if it's in the ConfigMap and the data entry key ends with `.tmpl` extension.

Here's an example of a ConfigMap containing a templated file for defining backends (`backends.vcl.tmpl`) and a static file that has the rest of the VCL configuration (`main.vcl`):

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: varnish-vcl-files
data:
  # A templated file as has the `.tmpl` extension
  backends.vcl.tmpl: |-
    import directors;

    #Backend nodes:
    {{ if .Backends -}}
    {{ range .Backends }}
    backend {{ .PodName }} {
      // backend {{ .PodName }} labels:
      {{- range $item, $key := .NodeLabels }}
      //   {{ $item }}: {{ $key -}}
      {{ end }}
      .host = "{{ .IP }}";
      .port = "{{ $.TargetPort }}";
    }
    {{ end }}
    {{- else -}}
    
    ...
  # A static file as doesn't has the `.tmpl` extension
  main.vcl: |
    vcl 4.0;

    import std;
    import var;
    import blob;
    include "backends.vcl";

    sub vcl_init {
      call init_backends;
      return (ok);
    }

    sub vcl_backend_response {
  ...
```

These are the available fields in the template that can be used to build your VCL files:

* `.Backends` - `[]PodInfo`: array of backends
  * `.IP` - `string`: IP address of a backend
  * `.NodeLabels` - `map[string]string`: labels of the node on which the backend is deployed.
  * `.PodName` - `string`: name of the pod representing a backend
* `.TargetPort` - `int`: port that is exposed on the backends
* `.VarnishNodes` - `[]PodInfo`: array of varnish nodes. Can be used for configuration of shard director (can be ignored if using a simple round robin director)
  * `.IP` - `string`: IP address of a varnish node
  * `.NodeLabels` - `map[string]string`: labels of the node on which a varnish node is deployed.
  * `.PodName` - `string`: name of the pod representing a varnish node
* `.VarnishPort` - `int`: port that is exposed on varnish nodes


For example, to generate your `backend`'s definitions you can use the following template:

```vcl
{{ range .Backends }}
backend {{ .PodName }} {
  .host = "{{ .IP }}";
  .port = "{{ $.TargetPort }}";
}
{{ end }}
```

This loops over the `.Backends` array, names each backend `.PodName`, sets `.host` to `.IP`, and then sets port to the universal `$.TargetPort`.

Assuming you have 3 nginx backends that listen on port 80, the resulting code can look like this:

```vcl
backend nginx-backend-6f4c6cbc6c-mjpck {
  .host = "172.30.243.4";
  .port = "80";
}

backend nginx-backend-6f4c6cbc6c-jmhhb {
  .host = "172.30.246.131";
  .port = "80";
}

backend nginx-backend-6f4c6cbc6c-ckqmv {
  .host = "172.30.80.137";
  .port = "80";
}
```

### Using User Defined VCL Code Versions

VCL related status information is available at field `.status.vcl`. 

The current VCL version can be found at `.status.vcl.configMapVersion`. It matches the resource version of the ConfigMap that contains the VCL code. 

To tag your own user friendly versions, an annotation `VCLVersion` on the ConfigMap can be used.

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

After setting the annotation, that version can be seen at `.status.vcl.version`. This field is optional and will only appear if the annotation is set.

```yaml
apiVersion: icm.ibm.com/v1alpha1
kind: VarnishService
metadata:
    ...
status:
  vcl:
    version: v1.0 # <-- reflects the `VCLVersion` annotation in the ConfigMap
    configMapVersion: "292181" # <-- reflects the ConfigMap resource version
  ...
```

After the VCL in the ConfigMap has been changed, the associated status fields will be immediately updated to reflect the latest version. However that does not guarantee that Varnish pods run the latest VCL configuration. It needs time to reload and if there is a problem, such as a syntax error in the VCL, Varnish will not load until the VCL is fixed.
 
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
  statefulSet:
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
NAME                              CONFIG_MAP_VERSION
my-varnish-varnish-statefulset-0  351231
my-varnish-varnish-statefulset-1  351231
my-varnish-varnish-statefulset-2  351231
# Looks like all pods have outdated VCL. Lets check the logs of one of the pods
> kubectl logs -n my-varnish my-varnish-varnish-statefulset-0 
2019-06-24T12:59:56.105Z	INFO	controller/controller_files.go:57	Rewriting file	{"kwatcher_version": "0.14.6", "varnish_service": "my-varnish", "pod_name": "my-varnish-varnish-statefulset-0", "namespace": "my-varnish", "file_path": "/etc/varnish/backends.vcl"}
2019-06-24T12:59:56.427Z	WARN	controller/controller_varnish.go:51	Message from VCC-compiler:
Expected one of
	'acl', 'sub', 'backend', 'probe', 'import', 'vcl',  or 'default'
Found: 'fsdfd' at
('/etc/varnish/backends.vcl' Line 8 Pos 1)
fsdfd
#####

Running VCC-compiler failed, exited with 2
Command failed with error code 106
VCL compilation failed
No VCL named v-20861922-1561381196 known.
Command failed with error code 106
	{"kwatcher_version": "0.14.6", "varnish_service": "my-varnish", "pod_name": "my-varnish-varnish-statefulset-0", "namespace": "my-varnish"}
```

As the logs indicate, the issue here is the invalid VCL syntax.