package controller

var haproxyConfigTemplate = `
global
  daemon
  maxconn {{ .MaxConnections }}
  log stdout format raw daemon

defaults
  log global
  mode http
  timeout connect {{ .ConnectTimeout }}ms
  timeout client {{ .ClientTimeout }}ms
  timeout server {{ .ServerTimeout }}ms

frontend stats
  bind *:8404
  option http-use-htx
  http-request use-service prometheus-exporter if { path /metrics }
  stats enable
  stats uri /stats
  stats refresh {{ .StatRefreshRate }}s

frontend localhost
  # Only bind on 80 if you also want to listen for connections on 80
  bind *:8080
  mode tcp
  option tcplog
  default_backend nodes

backend nodes
  mode http
  balance roundrobin
  http-response set-header Strict-Transport-Security max-age={{ .BackendServerMaxAgeHeader }}
  http-request set-header Host {{ .BackendServerHostHeader }}
  {{ range $i, $server := .BackendServers }}
  server svr{{ $i }} {{ $server }}:{{ $.BackendServerPort }} ssl sni str({{ $server }}) verify none
  {{ end }}
`
