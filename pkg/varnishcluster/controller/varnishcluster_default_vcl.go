package controller

// VCL files that are used to create a default ConfigMap.

const entrypointVCLFileContent = `vcl 4.0;

import std;
import var;
include "backends.vcl";

sub vcl_init {
  call init_backends;
}

sub vcl_recv {
  if (req.restarts > 0) {
    set req.hash_always_miss = true;
  }

  if (req.method == "GET" && req.url == "/heartbeat") {
    return(synth(200, "OK"));
  }

  // If backends are not configured correctly
  if (!(var.global_get("backendsFound") == "true")) {
    return(synth(503, "No backends configured"));
  }

  set req.backend_hint = container_rr.backend();

  if (req.method == "GET" && req.url == "/liveness") {
    if (!std.healthy(req.backend_hint)) {
      return(synth(503, "No healthy backends"));
    }
    return(synth(200, "OK"));
  }

  // Do not cache paths with health (healthcheck cache fix)
  if (req.url ~ "health") {
    return (pass);
  }
}

sub vcl_synth {
       set resp.http.Content-Type = "text/html; charset=utf-8";

       if (!(var.global_get("backendsFound") == "true")) { //error message if no backends configured
          synthetic( {"<!DOCTYPE html>
           <html>
             <head>
               <title>Incorrect backend configuration"</title>
             </head>
             <body>
               <h1>Incorrect backend configuration</h1>
               <p>Please check your deployment. It may not have pods running or Varnish is pointed to a non existing deployment.</p>
               <p>XID: "} + req.xid + {"</p>
               <hr>
             </body>
           </html>
           "} );
       } else { //default error message for the rest of the cases
        synthetic( {"<!DOCTYPE html>
            <html>
              <head>
                <title>"} + resp.status + " " + resp.reason + {"</title>
              </head>
              <body>
                <h1>Error "} + resp.status + " " + resp.reason + {"</h1>
                <p>"} + resp.reason + {"</p>
                <h3>Guru meditation:</h3>
                <p>XID: "} + req.xid + {"</p>
                <hr>
                <p>Varnish cache server</p>
              </body>
            </html>
            "} );
       }

    return (deliver);
}

sub vcl_hit {
  // Do not serve stale objects
  if (obj.ttl >= 0s) {
    return (deliver);
  }
  return (restart);
}

sub vcl_backend_response {

  // Do not cache 404s from backends
  if (beresp.status == 404) {
    set beresp.ttl = 0s;
  }
}

sub vcl_deliver {

  set resp.http.grace = req.http.grace;

  if (obj.hits > 0) {
    set resp.http.X-Varnish-Cache = "HIT";
  }
  else {
    set resp.http.X-Varnish-Cache = "MISS";
  }

  return (deliver);
}

`

const backendsVCLTmplFileContent = `import directors;

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
// Without this dummy backend, varnish will not compile the code
// This is a dummy, and should not be used anywhere
backend dummy {
  .host = "127.0.0.1";
  .port = "0";
}
{{- end }}

sub init_backends {
  // The line below is generated and creates a variable that is used to build custom logic
  // when the user configured the backends incorrectly. E.g. return a custom error page that indicates the issue.
  var.global_set("backendsFound", {{ if .Backends }}"true"{{ else }}"false"{{ end }}); //only strings are allowed to be set globally

  new container_rr = directors.round_robin();
  {{- range .Backends }}
  container_rr.add_backend({{ .PodName }});
  {{- end }}
}
`
