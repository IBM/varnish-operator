import directors;

// Without this dummy backend, varnish will not compile the code
// This is a dummy, and should not be used anywhere
backend dummy {
  .host = "127.0.0.1";
  .port = "0";
}

sub init_backends {
  new container_rr = directors.round_robin();
}