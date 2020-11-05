# VCL files configuration

Varnish VCL files are stored in a ConfigMap. Each entry in the ConfigMap corresponds to a file with VCL configuration. 

Modifications to VCL files in the ConfigMap are automatically applied to Varnish pods. This is handled by doing a VCL reload that does not restart the Varnish process. So the existing cache will not be lost.

There are 2 fields relevant to configuring the `VarnishCluster` for VCL code, in `.spec.vcl` object:

* **configMapName**: This is a required field and tells the `VarnishCluster` resource the name of the ConfigMap that contains/will contain the VCL files
* **entrypointFileName**: The name of the file that acts as the entrypoint for Varnish. If the entrypoint file is templated (ends in `.tmpl`), exclude the `.tmpl` extension. For example, if ConfigMap has file `mytemplatedfile.vcl.tmpl`, set `entrypointFileName: mytemplatedfile.vcl` as the generated file will omit the extension.

If a ConfigMap does not exist on `VarnishCluster` creation, the operator will create one and populate it with a default `backends.vcl.tmpl` and `entrypoint.vcl`. Their behavior is as follows:

* `backends.vcl.tmpl`: collect all backends into a single round-robin director
* `entrypoint.vcl`:
  * respond to `GET /heartbeat` checks with a 200
  * respond to `GET /liveness` checks with a 200 or 503, depending on healthy backends
  * respond to all other requests normally, caching all non-404 responses
  * hash request based on url
  * add `X-Varnish-Cache` header to response with "HIT" or "MISS" value, based on presence in cache

### Writing a Templated VCL File

The template file is a regular VCL file, with the addition of [Go templates](https://golang.org/pkg/text/template). This is because there is no way to know the backend's IP addresses at startup, so they must be injected at runtime. Also they can change over time if the backends get rescheduled by Kubernetes. 

The file is considered templated if it's in the ConfigMap and the data entry key ends with `.tmpl` extension.

Here's an example of a ConfigMap containing a templated file for defining backends (`backends.vcl.tmpl`) and a static file that has the rest of the VCL configuration (`entrypoint.vcl`):

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
  entrypoint.vcl: |
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
  * `.Weight` - `float64`: backend weight
  {% hint style="info" %}
  Please note that only the Random director can accept Weight as backend parameter
  [Random director documentation](https://varnish-cache.org/docs/6.1/reference/vmod_directors.generated.html?highlight=round%20robin#void-xrandom-add-backend-backend-real)
  For more information regarding weight control see [VarnishCluster](varnish-cluster.md)
  {% endhint %}
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

The `.Backends` array includes all backend nodes, no matter if the instance is ready or not. For that reason, it is strongly recommended to set [health checks](https://varnish-cache.org/docs/3.0/tutorial/advanced_backend_servers.html#health-checks) in your VCL configuration. Otherwise, Varnish could send traffic to not yet ready backends.

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
apiVersion: caching.ibm.com/v1alpha1
kind: VarnishCluster
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
apiVersion: caching.ibm.com/v1alpha1
kind: VarnishCluster
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

To check which pods have outdated versions, simply check their annotations. The annotation `configMapVersion` on the Varnish pod will indicate the latest version of the ConfigMap used. If it's not the same as in the `VarnishCluster` status it's likely that there's an issue.

Example of detecting a pod that failed to reload:

```bash
# get the latest version
> kubectl get varnishcluster -n varnish-ns my-varnish -o=custom-columns=NAME:.metadata.name,CONFIG_MAP_VERSION:.status.vcl.configMapVersion
NAME        CONFIG_MAP_VERSION
my-varnish  292181
# figure out which pods doesn't have that latest version
> kubectl get pods -n varnish-ns -o=custom-columns=NAME:.metadata.name,CONFIG_MAP_VERSION:.metadata.annotations.configMapVersion
NAME                              CONFIG_MAP_VERSION
my-varnish-varnish-0  351231
my-varnish-varnish-1  351231
my-varnish-varnish-2  351231
# Looks like all pods have outdated VCL. Lets check the logs of one of the pods
> kubectl logs -n my-varnish my-varnish-varnish-0 
2019-06-24T12:59:56.105Z	INFO	controller/controller_files.go:57	Rewriting file	{"varnish_controller_version": "0.21.0", "varnish_cluster": "my-varnish", "pod_name": "my-varnish-varnish-0", "namespace": "my-varnish", "file_path": "/etc/varnish/backends.vcl"}
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
	{"varnish_controller_version": "0.21.0", "varnish_service": "my-varnish", "pod_name": "my-varnish-varnish-0", "namespace": "my-varnish"}
```

As the logs indicate, the issue here is the invalid VCL syntax.

### Passing additional information into VCL

The `VarnishCluster` spec has a field `.spec.varnish.envFrom` that allows injecting custom values into env vars. After defining, in VCL files you can read them using [std.getenv()](https://varnish-cache.org/docs/5.1/reference/vmod_std.generated.html#func-getenv) function. Bot Secrets and ConfigMaps can be used to do it.

Example of injecting values from a Secret:

1. Create a secret: `kubectl create secret generic vcl-secrets --from-literal=secret-cookie=secret-cookie-value` 

2. Use it by specifying it in the `VarnishCluster` spec

```yaml
apiVersion: caching.ibm.com/v1alpha1
kind: VarnishCluster
metadata:
  name: example
spec:
  varnish:
    envFrom:
    - secretRef:
        name: vcl-secrets
```

3. Read the injected env var in you VCL code

```c
   ...
    sub vcl_recv {
      // validate cookie
      if (!req.http.Cookie ~ std.getenv("secret-cookie")) {
          return (synth(401));
      }
   ... 
```

That way, your secret value can be stored a separately from VCL and used securely.

Note that since those values passed using env vars, if a value changes, a reload of a pod needed to pass the new value to Varnish. Which means cache loss.

### Advanced VCL with Varnish sharding

Starting from Varnish 5.0 directors VMOD were extended by adding new director - shard director. It's basically a better hash director. The shard director selects backends by a key, which can be provided directly or derived from strings. For the same key, the shard director will always return the same backend, unless the backend configuration or health state changes. For differing keys, the shard director will likely choose different backends. In the default configuration, unhealthy backends are not selected.
When used with clusters of Varnish servers, the shard director will, if otherwise configured equally, make the same decision on all servers. So requests sharing a common criterion used as the shard key will be balanced onto the same backend servers no matter which Varnish server handles the request.

There are couple of very useful options that can be passed to shard director. The `rampup` feature is a slow start mechanism that allows just-gone-healthy backends to ease into full load smoothly, while the `warmup` feature prepares backends for the traffic they would see if the primary backend for a key goes down.

#### Creating shard director

Leveraging templating capabilities provided by operator, we can build `backends.vcl` config file. It will contain all entities required to setup varnish sharding cluster:
 - Application backends
 - Varnish backends (cluster members)
 - Heartbeat probe for Varnish instances
 - Director that holds application backends (round-robin, random, etc.)
 - Shard director
 - ACL with Varnish cluster members

As described in [Writing a Templated VCL File](#writing-a-templated-vcl-file) section, create `backends.vcl.tmpl` and create a structure with all components:

```vcl
// Import VMOD directors
import directors;

// Define probe used for heartbeat
probe heartbeat {
  .request = "HEAD /heartbeat HTTP/1.1"
      "Connection: close"
      "Host: shard";
  .interval = 1s;
}

// Application backends
{{ range .Backends }}
backend {{ .PodName }} {
  .host = "{{ .IP }}";
  .port = "{{ $.TargetPort }}";
}
{{ end }}

// Varnish cluster backends
{{ range .VarnishNodes }}
backend {{ .PodName }} {
  .host = "{{ .IP }}";
  .port = "{{ $.VarnishPort }}";
  .probe = heartbeat;
}
{{ end }}

// Create ACL with Varnish cluster members
acl acl_cluster {
  {{ range .VarnishNodes }}
  "{{ .IP }}"/32;
  {{ end }}
}

// Here we create two directors - "real", application round-robin
// and "cluster", Varnish shard director
sub init_backends {
  new real = directors.round_robin();
  {{- range .Backends }}
  real.add_backend({{ .PodName }});
  {{- end }}

  new cluster = directors.shard();
  {{ range .VarnishNodes }}
  cluster.add_backend({{ .PodName }});
  {{ end }}
  cluster.set_rampup(30s);
  cluster.set_warmup(0.1);
  cluster.reconfigure();
}
```
Please note `set_rampup` and `set_warmup` options being passed to `cluster` director. They provide a mechanism to change how the shard director will manage cluster members downtime. 

`set_rampup` configures slow start interval in seconds for servers which just came back from unhealthy. If chosen backend is in its rampup period, with a probability proportional to the fraction of time since the backup became healthy to the rampup period, return the next alternative backend, unless this is also in its rampup period. If duration is `0` (default), rampup is disabled. 

`set_warmup` configures the default warmup probability. Sets the ratio of requests (`0.0` to `1.0`) that goes to the next alternate backend to warm it up when the preferred backend is healthy. Not active if any of the preferred or alternative backend are in rampup. Setting of `0.5` is a convenient way to spread the load for each key over two backends under normal operating conditions. If probability is `0.0` (default), warmup is disabled.

Note that the shard director needs to be configured using at least one `.add_backend()` call(s) followed by a `.reconfigure()` call before it can hand out backends.

Now, once backends and directors are properly configured, we can proceed with configuring main logic and client requests are processed. Next configuration snippet is **not complete VCL**, but rather only client facing logic - `vcl_recv` subroutine. 

```vcl
// Include and initialize backends
include "backends.vcl";

sub vcl_init {
  call init_backends;
  return (ok);
}

sub vcl_recv {

  // Answer to heartbeats only from cluster members
  if (remote.ip ~ acl_cluster) {
    if (req.http.Host == "shard") {
      if (req.url == "/heartbeat") {
        return (synth(200));
      }
      return (synth(404));
    }
  }

  // Let shard director to pick the backend (shard)
  set req.backend_hint = cluster.backend(URL);
  set req.http.X-shard = req.backend_hint;

  // Use application backend if request came here from ourselves or from
  // cluster members (i.e. sharding is already happened), otherwise let it pass
  // to another shard
  if (req.http.X-shard == server.identity || remote.ip ~ acl_cluster) {
    set req.backend_hint = real.backend();
  } else {
    return(pass);
  }

  return (hash);
}
...
```
This snippet will route all requests to shard director for sharding, prevent multiple in-cluster hops when warmup is happened, and reroute requests to the shard that it primary for given request.

Everything together can provide reliable way of handling lifecycle events for Varnish cluster members with graceful fallback and controlled pre-warmup of alternative backends.

