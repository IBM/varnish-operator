vcl 4.0;

import std;
include "backends.vcl";

sub vcl_init {
  call init_backends;
  return (ok);
}

sub vcl_recv {

  if (req.method == "GET" && req.url == "/heartbeat") {
    return(synth(200, "OK"));
  }
  
  set req.backend_hint = container_rr.backend();

  if (req.method == "GET" && req.url == "/liveness") {
    if (!std.healthy(req.backend_hint)) {
      return(synth(503, "No healthy backends"));
    }
    return(synth(200, "OK"));
  }

  #
  # Do not cache paths with health (healthcheck cache fix)
  ####
  if (req.url ~ "health") {
    return (pass);
  }

  return (hash);
}

sub vcl_hash {

  # Called after vcl_recv to create a hash value for the request. This is used as a key
  # to look up the object in Varnish.
  hash_data(req.url);

  return (lookup);
}

sub vcl_hit {
  # Do not serve stale objects
  if (obj.ttl >= 0s) {
    return (deliver);
  }
  return (miss);
}

sub vcl_backend_response {

  # Do not cache 404s from backends
  if (beresp.status == 404) {
    set beresp.ttl = 0s;
  }

  return (deliver);
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

sub vcl_fini {

  return (ok);
}