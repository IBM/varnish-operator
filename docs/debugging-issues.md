# Debugging issues

## Log levels
The operator has configurable log levels which can be set in the `logLevel` field of your `values.yaml` [Helm chart](operator-configuration.md) override. Also you can change the output of the logs using `logFormat` field. It is `json` by default, but you can set it to `console` to have more user friendly log output.

The `VarnishService` has the same configuration options regarding logging with `spec.logLevel` and `spec.logFormat` fields in the [`VarnishService` spec](varnish-service-configuration.md). Keep in mind that changing log level will cause the pods to restart and invalidate the cache.

## Using Varnish tools

To debug some Varnish related issues you may want to use the tools provided by Varnish (`varnishlog`, `varnishadm`, `varnishncsa`, etc.). Those tools are available in the containers Varnish is running in.

After you've [created your `VarnishService`](varnish-service.md) you should be able to see your Varnish pods. You can use the `varnish-owner=<your-varnishservice-name>` label to select your pods.

For a `VarnishService` named `varnish-service-example` the command will look like this:

```bash
$ kubectl get pods -l varnish-owner=varnish-service-example                                                   
NAME                                                          READY   STATUS    RESTARTS   AGE
varnish-service-example-varnish-deployment-6875b997cf-mvxp5   1/1     Running   0          15s
```

Now, you can exec to that pod's container:

```bash
$ kubectl exec -it varnish-service-example-varnish-deployment-6875b997cf-mvxp5 sh
```

and execute any available Varnish command line tool. For example `varnishadm`:

```bash
$ varnishadm param.show default_ttl # find out the default TTL
default_ttl
        Value is: 120.000 [seconds] (default)
        Minimum is: 0.000

        The TTL assigned to objects if neither the backend nor the VCL
        code assigns one.

        NB: This parameter is evaluated only when objects are created.
        To change it for all objects, restart or ban everything.
```

or check Varnish logs (need to make requests to Varnish to actually see some logs):

```bash
$ varnishlog # send some requests to make logs appear
*   << Request  >> 32770     
-   Begin          req 32769 rxreq
-   Timestamp      Start: 1562234840.724019 0.000000 0.000000
-   Timestamp      Req: 1562234840.724019 0.000000 0.000000
-   ReqStart       127.0.0.1 57828 a0
-   ReqMethod      GET
-   ReqURL         /
-   ReqProtocol    HTTP/1.1
-   ReqHeader      Host: localhost:8080
-   ReqHeader      Connection: keep-alive
-   ReqHeader      Upgrade-Insecure-Requests: 1
-   ReqHeader      User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/75.0.3770.100 Safari/537.36
-   ReqHeader      Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3
-   ReqHeader      Accept-Encoding: gzip, deflate, br
-   ReqHeader      Accept-Language: en-US,en;q=0.9
-   ReqHeader      If-None-Match: "5ce409fd-264"
-   ReqHeader      If-Modified-Since: Tue, 21 May 2019 14:23:57 GMT
-   ReqHeader      X-Forwarded-For: 127.0.0.1
-   VCL_call       RECV
-   VCL_return     synth
-   ReqUnset       Accept-Encoding: gzip, deflate, br
-   ReqHeader      Accept-Encoding: gzip
-   VCL_call       HASH
-   VCL_return     lookup
-   Timestamp      Process: 1562234840.724098 0.000079 0.000079
-   RespHeader     Date: Thu, 04 Jul 2019 10:07: 
....
```

Also you can use `varnishncsa` tool to debug requests. For example, to see HIT/MISS for requests in real time:

```bash
$ varnishncsa -F '%m %U%q - %{Varnish:hitmiss}x' # send requests to see some output
GET / - miss
GET /favicon.ico - miss
GET / - hit
GET / - hit
GET / - hit
GET /favicon.ico - miss

```

Note that if you have multiple replicas, Kubernetes will balance the load between them and you may not see a log entry as the request could have gone to a different replica. Just make a few more requests and one of them will eventually land on the pod you're monitoring.

### VarnishService status

As with any other Kubernetes resource, `VarnishService` has a `status` object describing the state of the resource. You can find the status of underlying Deployment and Service objects in a similar manner. The `status` object can also reveal information about the status of Varnish VCL configuration. It can be found in the `.status.vcl` object. The `status.vcl.availability` field is especially useful for debugging. It shows how many Varnish instances are running with the latest version of VCL.

For example, if the field has value `availability: 3 latest / 0 outdated` it means that all pods are running the latest version of VCL. However if none of the pods are running the latest version (`availability: "0 latest / 3 outdated"`), it could mean that the VCL could be invalid. You can check it by looking at `VarnishService` events first:

```bash
$ kubectl describe vs varnish-service-example
Name:         varnish-service-example
Namespace:    varnish-operator-system
    ...
  Vcl:
    Availability:        0 latest / 1 outdated
    Config Map Version:  974330
Events:
  Type     Reason               Age   From     Message
  ----     ------               ----  ----     -------
  Warning  VCLCompilationError  11s   varnish  VCL compilation failed for pod varnish-service-example-varnish-deployment-6875b997cf-mvxp5. See pod logs for details
```

As you can see, the event indicates that it is a VCL compilation error indeed. To check the compilation error message see the pod logs:

```bash
$ kubectl logs varnish-service-example-varnish-deployment-5c84d4c876-45qrh                 
...
2019-07-04T10:50:53.481Z	INFO	kwatcher/main.go:60	Starting Varnish Watcher	{"kwatcher_version": "0.17.0"}
2019-07-04T10:50:54.012Z	INFO	controller/controller_files.go:57	Rewriting file	{"kwatcher_version": "0.17.0", "varnish_service": "varnish-service-example", "pod_name": "varnish-service-example-varnish-deployment-5c84d4c876-45qrh", "namespace": "varnish-operator-system", "file_path": "/etc/varnish/backends.vcl"}
2019-07-04T10:50:54.506Z	WARN	controller/controller_varnish.go:51	Message from VCC-compiler:
Expected one of
	'acl', 'sub', 'backend', 'probe', 'import', 'vcl',  or 'default'
Found: 'bakcend' at
('/etc/varnish/backends.vcl' Line 6 Pos 1)
bakcend nginx-5c7588df-rcrn4 {
#######-----------------------

Running VCC-compiler failed, exited with 2
Command failed with error code 106
VCL compilation failed
No VCL named v-974330-1562237454 known.
Command failed with error code 106
	{"kwatcher_version": "0.17.0", "varnish_service": "varnish-service-example", "pod_name": "varnish-service-example-varnish-deployment-5c84d4c876-45qrh", "namespace": "varnish-operator-system"}
```

As you can see we have a type in word - `bakcend`. To fix it you'll have to modify your [ConfigMap containing your VCL files](vcl-configuration.md).
