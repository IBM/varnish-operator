package controller

var haproxyGrafanaDashboardTemplate = `
{
  "__requires": [
    {
      "type": "grafana",
      "id": "grafana",
      "name": "Grafana",
      "version": "7.3.7"
    },
    {
      "type": "panel",
      "id": "graph",
      "name": "Graph",
      "version": ""
    },
    {
      "type": "datasource",
      "id": "prometheus",
      "name": "Prometheus",
      "version": "1.0.0"
    },
    {
      "type": "panel",
      "id": "singlestat",
      "name": "Singlestat",
      "version": ""
    }
  ],
  "annotations": {
    "list": [
      {
        "$$hashKey": "object:257",
        "builtIn": 1,
        "datasource": "-- Grafana --",
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "description": "HAProxy with Prometheus data",
  "editable": true,
  "gnetId": 12693,
  "graphTooltip": 1,
  "id": null,
  "iteration": 1621339804258,
  "links": [],
  "panels": [
    {
      "collapsed": false,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 0
      },
      "id": 152,
      "panels": [],
      "repeat": null,
      "title": "Basic General Info",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "{{.DatasourceName}}",
      "description": "",
      "editable": true,
      "error": false,
      "fieldConfig": {
        "defaults": {
          "custom": {},
          "links": []
        },
        "overrides": []
      },
      "fill": 2,
      "fillGradient": 0,
      "grid": {},
      "gridPos": {
        "h": 10,
        "w": 24,
        "x": 0,
        "y": 1
      },
      "hiddenSeries": false,
      "id": 83,
      "legend": {
        "alignAsTable": true,
        "avg": true,
        "current": true,
        "max": true,
        "min": true,
        "rightSide": true,
        "show": true,
        "total": false,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "maxPerRow": 2,
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.3.7",
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "$$hashKey": "object:469",
          "alias": "/.*Back.*/",
          "transform": "negative-Y"
        },
        {
          "$$hashKey": "object:470",
          "alias": "/.*1.*/",
          "color": "#6ED0E0"
        },
        {
          "$$hashKey": "object:471",
          "alias": "/.*2.*/",
          "color": "#7EB26D"
        },
        {
          "$$hashKey": "object:472",
          "alias": "/.*3.*/",
          "color": "#1F78C1"
        },
        {
          "$$hashKey": "object:473",
          "alias": "/.*4.*/",
          "color": "#CCA300"
        },
        {
          "$$hashKey": "object:474",
          "alias": "/.*5.*/",
          "color": "#890F02"
        },
        {
          "$$hashKey": "object:475",
          "alias": "/.*other.*/",
          "color": "#806EB7"
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(rate(haproxy_frontend_http_responses_total{proxy=~\"$frontend\",code=~\"$code\",instance=\"$host\"}[$__rate_interval])) by (code)",
          "interval": "$interval",
          "intervalFactor": 1,
          "legendFormat": "Front {{"{{ code }}"}}",
          "metric": "",
          "refId": "A",
          "step": 240
        },
        {
          "expr": "sum(rate(haproxy_backend_http_responses_total{proxy=~\"$backend\",code=~\"$code\",instance=\"$host\"}[$__rate_interval])) by (code)",
          "interval": "$interval",
          "intervalFactor": 1,
          "legendFormat": "Back {{"{{ code }}"}}",
          "metric": "",
          "refId": "B",
          "step": 240
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "All HTTP responses",
      "tooltip": {
        "msResolution": true,
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "$$hashKey": "object:524",
          "format": "short",
          "label": "- back / + front",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "$$hashKey": "object:525",
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": 0,
          "show": false
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "{{.DatasourceName}}",
      "decimals": 1,
      "description": "",
      "editable": true,
      "error": false,
      "fieldConfig": {
        "defaults": {
          "custom": {},
          "links": []
        },
        "overrides": []
      },
      "fill": 2,
      "fillGradient": 0,
      "grid": {},
      "gridPos": {
        "h": 10,
        "w": 12,
        "x": 0,
        "y": 11
      },
      "hiddenSeries": false,
      "id": 75,
      "legend": {
        "alignAsTable": true,
        "avg": true,
        "current": true,
        "max": true,
        "min": true,
        "rightSide": false,
        "show": true,
        "sideWidth": null,
        "total": false,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "maxPerRow": 2,
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.3.7",
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "$$hashKey": "object:2175",
          "alias": "/.*OUT.*/",
          "transform": "negative-Y"
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(rate(haproxy_frontend_bytes_in_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])*8) by (instance)",
          "interval": "$interval",
          "intervalFactor": 1,
          "legendFormat": "IN Front",
          "metric": "",
          "refId": "A",
          "step": 240
        },
        {
          "expr": "sum(rate(haproxy_frontend_bytes_out_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])*8) by (instance)",
          "interval": "$interval",
          "intervalFactor": 2,
          "legendFormat": "OUT Front",
          "refId": "B",
          "step": 240
        },
        {
          "expr": "sum(rate(haproxy_backend_bytes_in_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])*8) by (instance)",
          "intervalFactor": 2,
          "legendFormat": "IN Back",
          "refId": "C",
          "step": 240
        },
        {
          "expr": "sum(rate(haproxy_backend_bytes_out_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])*8) by (instance)",
          "intervalFactor": 2,
          "legendFormat": "OUT Back",
          "refId": "D",
          "step": 240
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Incoming / Outgoing bytes",
      "tooltip": {
        "msResolution": true,
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "$$hashKey": "object:2188",
          "format": "bits",
          "label": "- out / + in",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "$$hashKey": "object:2189",
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": 0,
          "show": false
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "{{.DatasourceName}}",
      "description": "",
      "editable": true,
      "error": false,
      "fieldConfig": {
        "defaults": {
          "custom": {},
          "links": []
        },
        "overrides": []
      },
      "fill": 2,
      "fillGradient": 0,
      "grid": {},
      "gridPos": {
        "h": 10,
        "w": 12,
        "x": 12,
        "y": 11
      },
      "hiddenSeries": false,
      "id": 79,
      "legend": {
        "alignAsTable": true,
        "avg": true,
        "current": true,
        "max": true,
        "min": true,
        "rightSide": false,
        "show": true,
        "total": false,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "maxPerRow": 2,
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.3.7",
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "$$hashKey": "object:616",
          "alias": "/.*Back.*/",
          "color": "#F2495C",
          "transform": "negative-Y"
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(rate(haproxy_frontend_connections_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])) by (instance)",
          "interval": "$interval",
          "intervalFactor": 1,
          "legendFormat": "Front",
          "metric": "",
          "refId": "A",
          "step": 240
        },
        {
          "expr": "sum(rate(haproxy_backend_connection_errors_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])) by (instance)",
          "hide": false,
          "interval": "$interval",
          "intervalFactor": 1,
          "legendFormat": "Back errors",
          "metric": "",
          "refId": "C",
          "step": 240
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Number of connections",
      "tooltip": {
        "msResolution": true,
        "shared": false,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "$$hashKey": "object:629",
          "format": "short",
          "label": "- back / + front",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "$$hashKey": "object:630",
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": 0,
          "show": false
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "{{.DatasourceName}}",
      "description": "",
      "editable": true,
      "error": false,
      "fieldConfig": {
        "defaults": {
          "custom": {},
          "links": []
        },
        "overrides": []
      },
      "fill": 2,
      "fillGradient": 0,
      "grid": {},
      "gridPos": {
        "h": 10,
        "w": 12,
        "x": 0,
        "y": 21
      },
      "hiddenSeries": false,
      "id": 81,
      "legend": {
        "alignAsTable": true,
        "avg": true,
        "current": true,
        "max": true,
        "min": true,
        "rightSide": false,
        "show": true,
        "total": false,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "maxPerRow": 2,
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.3.7",
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "$$hashKey": "object:931",
          "alias": "/.*Back.*/",
          "transform": "negative-Y"
        },
        {
          "$$hashKey": "object:1328",
          "alias": "/.*errors.*/",
          "color": "#F2495C"
        },
        {
          "$$hashKey": "object:1414",
          "alias": "/.*warn.*/",
          "color": "#FF9830"
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(rate(haproxy_frontend_http_requests_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])) by (instance)",
          "interval": "$interval",
          "intervalFactor": 1,
          "legendFormat": "Front requests",
          "metric": "",
          "refId": "A",
          "step": 240
        },
        {
          "expr": "sum(rate(haproxy_frontend_request_errors_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])) by (instance)",
          "interval": "$interval",
          "intervalFactor": 1,
          "legendFormat": "Front requests errors",
          "metric": "",
          "refId": "C",
          "step": 240
        },
        {
          "expr": "sum(rate(haproxy_frontend_requests_denied_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])) by (instance)",
          "interval": "$interval",
          "intervalFactor": 2,
          "legendFormat": "Front request denied",
          "refId": "F",
          "step": 240
        },
        {
          "expr": "sum(rate(haproxy_backend_redispatch_warnings_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])) by (instance)",
          "interval": "$interval",
          "intervalFactor": 2,
          "legendFormat": "Back redispatch warnings",
          "refId": "D",
          "step": 240
        },
        {
          "expr": "sum(rate(haproxy_backend_retry_warnings_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])) by (instance)",
          "interval": "$interval",
          "intervalFactor": 2,
          "legendFormat": "Back retry warnings",
          "refId": "E",
          "step": 240
        },
        {
          "expr": "sum(rate(haproxy_backend_response_errors_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])) by (instance)",
          "interval": "$interval",
          "intervalFactor": 1,
          "legendFormat": "Back response errors",
          "metric": "",
          "refId": "I",
          "step": 240
        },
        {
          "expr": "sum(haproxy_backend_current_queue{proxy=~\"$backend\",instance=\"$host\"}) by (instance)",
          "interval": "$interval",
          "intervalFactor": 2,
          "legendFormat": "Back queued requests",
          "refId": "G",
          "step": 240
        },
        {
          "expr": "sum(rate(haproxy_backend_http_requests_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])) by (instance)",
          "interval": "$interval",
          "intervalFactor": 1,
          "legendFormat": "Back requests",
          "metric": "",
          "refId": "H",
          "step": 240
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Requests and Responses",
      "tooltip": {
        "msResolution": true,
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "$$hashKey": "object:950",
          "format": "short",
          "label": "- back / + front",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "$$hashKey": "object:951",
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": 0,
          "show": false
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "{{.DatasourceName}}",
      "description": "",
      "editable": true,
      "error": false,
      "fieldConfig": {
        "defaults": {
          "custom": {},
          "links": []
        },
        "overrides": []
      },
      "fill": 2,
      "fillGradient": 0,
      "grid": {},
      "gridPos": {
        "h": 10,
        "w": 12,
        "x": 12,
        "y": 21
      },
      "hiddenSeries": false,
      "id": 84,
      "legend": {
        "alignAsTable": true,
        "avg": true,
        "current": true,
        "max": true,
        "min": true,
        "rightSide": false,
        "show": true,
        "total": false,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "maxPerRow": 2,
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.3.7",
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "$$hashKey": "object:1146",
          "alias": "/.*Back.*/",
          "transform": "negative-Y"
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(rate(haproxy_frontend_current_sessions{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])) by (instance)",
          "interval": "$interval",
          "intervalFactor": 1,
          "legendFormat": "Front",
          "metric": "",
          "refId": "B",
          "step": 240
        },
        {
          "expr": "sum(rate(haproxy_backend_current_sessions{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])) by (instance)",
          "interval": "$interval",
          "intervalFactor": 1,
          "legendFormat": "Back",
          "metric": "",
          "refId": "A",
          "step": 240
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Active sessions",
      "tooltip": {
        "msResolution": true,
        "shared": true,
        "sort": 2,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "$$hashKey": "object:1159",
          "format": "short",
          "label": "- back / + front",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "$$hashKey": "object:1160",
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": 0,
          "show": false
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "collapsed": false,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 31
      },
      "id": 151,
      "panels": [],
      "repeat": null,
      "title": "Basic General Status",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "{{.DatasourceName}}",
      "decimals": 0,
      "fieldConfig": {
        "defaults": {
          "custom": {},
          "links": []
        },
        "overrides": []
      },
      "fill": 5,
      "fillGradient": 0,
      "gridPos": {
        "h": 4,
        "w": 22,
        "x": 0,
        "y": 32
      },
      "hiddenSeries": false,
      "id": 85,
      "legend": {
        "alignAsTable": true,
        "avg": true,
        "current": true,
        "hideEmpty": false,
        "hideZero": false,
        "max": true,
        "min": true,
        "rightSide": true,
        "show": true,
        "total": false,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.3.7",
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "$$hashKey": "object:138",
          "alias": "Back Up",
          "transform": "negative-Y"
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "count(haproxy_frontend_status{instance=\"$host\"} == 1)",
          "hide": false,
          "interval": "$interval",
          "intervalFactor": 2,
          "legendFormat": "Front Up",
          "refId": "A",
          "step": 240
        },
        {
          "expr": "count(haproxy_backend_status{instance=\"$host\"} ==1)",
          "interval": "$interval",
          "intervalFactor": 2,
          "legendFormat": "Back Up",
          "refId": "B",
          "step": 240
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "All Status",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "$$hashKey": "object:155",
          "format": "short",
          "label": "- back / + front",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "$$hashKey": "object:156",
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": false
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "cacheTimeout": null,
      "colorBackground": false,
      "colorValue": false,
      "colors": [
        "rgba(245, 54, 54, 0.9)",
        "rgba(237, 129, 40, 0.89)",
        "rgba(50, 172, 45, 0.97)"
      ],
      "datasource": "{{.DatasourceName}}",
      "decimals": 0,
      "fieldConfig": {
        "defaults": {
          "custom": {}
        },
        "overrides": []
      },
      "format": "s",
      "gauge": {
        "maxValue": 100,
        "minValue": 0,
        "show": false,
        "thresholdLabels": false,
        "thresholdMarkers": true
      },
      "gridPos": {
        "h": 4,
        "w": 2,
        "x": 22,
        "y": 32
      },
      "id": 149,
      "interval": null,
      "links": [],
      "mappingType": 1,
      "mappingTypes": [
        {
          "$$hashKey": "object:71",
          "name": "value to text",
          "value": 1
        },
        {
          "$$hashKey": "object:72",
          "name": "range to text",
          "value": 2
        }
      ],
      "maxDataPoints": 100,
      "nullPointMode": "connected",
      "nullText": null,
      "postfix": " ago",
      "postfixFontSize": "50%",
      "prefix": "",
      "prefixFontSize": "50%",
      "rangeMaps": [
        {
          "from": "null",
          "text": "N/A",
          "to": "null"
        }
      ],
      "sparkline": {
        "fillColor": "rgba(31, 118, 189, 0.18)",
        "full": false,
        "lineColor": "rgb(31, 120, 193)",
        "show": false
      },
      "tableColumn": "",
      "targets": [
        {
          "expr": "time() - haproxy_process_start_time_seconds{instance=\"$host\"}",
          "intervalFactor": 2,
          "refId": "A",
          "step": 240
        }
      ],
      "thresholds": "",
      "title": "Started...",
      "type": "singlestat",
      "valueFontSize": "50%",
      "valueMaps": [
        {
          "$$hashKey": "object:74",
          "op": "=",
          "text": "N/A",
          "value": "null"
        }
      ],
      "valueName": "current"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 36
      },
      "id": 182,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 1,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 15,
            "w": 12,
            "x": 0,
            "y": 3
          },
          "hiddenSeries": false,
          "id": 42,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "sideWidth": null,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:3222",
              "alias": "/.*OUT.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_frontend_bytes_in_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])*8",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "IN {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_frontend_bytes_out_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])*8",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "OUT {{"{{ proxy }}"}}",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Incoming / Outgoing bytes",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3235",
              "format": "bits",
              "label": "- out / + in",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:3236",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 1,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 15,
            "w": 12,
            "x": 12,
            "y": 3
          },
          "hiddenSeries": false,
          "id": 1,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:3119",
              "alias": "/.*OUT.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_bytes_in_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])*8",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "IN {{"{{ proxy }}"}}",
              "metric": "haproxy_backend_",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_backend_bytes_out_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])*8",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "OUT {{"{{ proxy }}"}}",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Incoming / Outgoing bytes",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3132",
              "format": "bits",
              "label": "- out / + in",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:3133",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 15,
            "w": 12,
            "x": 0,
            "y": 18
          },
          "hiddenSeries": false,
          "id": 43,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:3369",
              "alias": "/.*Denied*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_frontend_connections_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Successful {{"{{ proxy }}"}}",
              "metric": "haproxy_backe",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_frontend_denied_connections_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Denied {{"{{ proxy }}"}}",
              "metric": "haproxy_backe",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Connections successful / denied",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3389",
              "format": "short",
              "label": "- denied / + successful",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:3390",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 15,
            "w": 12,
            "x": 12,
            "y": 18
          },
          "hiddenSeries": false,
          "id": 27,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:2890",
              "alias": "/.*Error.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_connection_attempts_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Attempts {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_backend_connection_errors_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Error {{"{{ proxy }}"}}",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Connections attempts / errors",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:2903",
              "format": "short",
              "label": "- error / + attempt",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:2904",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 15,
            "w": 12,
            "x": 0,
            "y": 33
          },
          "hiddenSeries": false,
          "id": 114,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_frontend_connections_rate_max{proxy=~\"$frontend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Max {{"{{ proxy }}"}}",
              "metric": "haproxy_backe",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Maximum observed number of connections per second",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:437",
              "format": "short",
              "label": "connections",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:438",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 15,
            "w": 12,
            "x": 12,
            "y": 33
          },
          "hiddenSeries": false,
          "id": 131,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_connection_reuses_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Reuses {{"{{ proxy }}"}}",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Connections reuses",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:356",
              "format": "short",
              "label": "reuses",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:357",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "title": "Throughtput / Connections",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 37
      },
      "id": 154,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 7
          },
          "hiddenSeries": false,
          "id": 28,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_current_queue{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Queued {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Queued requests not assigned to any server",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:781",
              "format": "short",
              "label": "requests",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:782",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 7
          },
          "hiddenSeries": false,
          "id": 32,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_max_queue{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Max {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Max queued requests not assigned to any server",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1157",
              "format": "short",
              "label": "requests",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:1158",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "repeat": null,
      "title": "Queues",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 38
      },
      "id": 155,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 0,
            "y": 5
          },
          "hiddenSeries": false,
          "id": 134,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:1396",
              "alias": "/.*Denied.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_http_requests_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_backend_requests_denied_total{proxy=~\"$frontend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Denied {{"{{ proxy }}"}}",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - HTTP requests OK / Denied",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1415",
              "format": "short",
              "label": "- denied / + ok",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:1416",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 12,
            "y": 5
          },
          "hiddenSeries": false,
          "id": 46,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:1641",
              "alias": "/.*Error.*/",
              "transform": "negative-Y"
            },
            {
              "$$hashKey": "object:1642",
              "alias": "/.*Denied.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_frontend_http_requests_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_frontend_request_errors_total{proxy=~\"$frontend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Error {{"{{ proxy }}"}}",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_frontend_requests_denied_total{proxy=~\"$frontend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Denied {{"{{ proxy }}"}}",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - HTTP requests OK / Error / Denied",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1661",
              "format": "short",
              "label": "- error - denied / + ok",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:1662",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 0,
            "y": 19
          },
          "hiddenSeries": false,
          "id": 126,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "hideEmpty": false,
            "hideZero": false,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_responses_denied_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Denied {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - HTTP responses denied",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1906",
              "format": "short",
              "label": "denied",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:1907",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 12,
            "y": 19
          },
          "hiddenSeries": false,
          "id": 115,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:2055",
              "alias": "/.*Denied.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_frontend_responses_denied_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Denied {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - HTTP responses denied",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:2074",
              "format": "short",
              "label": "denied",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:2075",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 0,
            "y": 33
          },
          "hiddenSeries": false,
          "id": 35,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "hideEmpty": false,
            "hideZero": false,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:2248",
              "alias": "/.*Error.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_redispatch_warnings_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Redispatch {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_backend_retry_warnings_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Retry {{"{{ proxy }}"}}",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_backend_response_errors_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Error {{"{{ proxy }}"}}",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Error / Redispatch / Retry",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:2261",
              "format": "short",
              "label": "- error / + redispatch + retry",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:2262",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 12,
            "y": 33
          },
          "hiddenSeries": false,
          "id": 138,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "hideEmpty": false,
            "hideZero": false,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:2857",
              "alias": "/.*.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_frontend_http_requests_rate_max{proxy=~\"$backend\", instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Max {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Maximum observed number of HTTP requests per second",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:2580",
              "format": "short",
              "label": "requests",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:2581",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "repeat": null,
      "title": "Requests / Responses",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 39
      },
      "id": 176,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 0,
            "y": 6
          },
          "hiddenSeries": false,
          "id": 132,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_connect_time_average_seconds{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "haproxy_backend_current_queue",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Avg connection time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3021",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:3022",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 12,
            "y": 6
          },
          "hiddenSeries": false,
          "id": 209,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_max_connect_time_seconds{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "haproxy_backend_current_queue",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Max connection time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3021",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:3022",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 0,
            "y": 20
          },
          "hiddenSeries": false,
          "id": 178,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_total_time_average_seconds{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "haproxy_backend_current_queue",
              "refId": "D",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Avg. total time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3021",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:3022",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 12,
            "y": 20
          },
          "hiddenSeries": false,
          "id": 210,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_max_total_time_seconds{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "haproxy_backend_current_queue",
              "refId": "D",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Max total time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3021",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:3022",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 0,
            "y": 34
          },
          "hiddenSeries": false,
          "id": 177,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_response_time_average_seconds{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "haproxy_backend_current_queue",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Avg. response time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3021",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:3022",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 12,
            "y": 34
          },
          "hiddenSeries": false,
          "id": 211,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_max_response_time_seconds{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "haproxy_backend_current_queue",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Max response time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3021",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:3022",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 0,
            "y": 48
          },
          "hiddenSeries": false,
          "id": 127,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_queue_time_average_seconds{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Avg. queue time for last 1024 successful connections",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1225",
              "decimals": null,
              "format": "s",
              "label": "",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:1226",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 12,
            "y": 48
          },
          "hiddenSeries": false,
          "id": 212,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_max_queue_time_seconds{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Max queue time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1225",
              "decimals": null,
              "format": "s",
              "label": "",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:1226",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "title": "Times",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 40
      },
      "id": 156,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 13
          },
          "hiddenSeries": false,
          "id": 47,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_frontend_http_responses_total{proxy=~\"$frontend\", code=~\"$code\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ code }}"}} {{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - HTTP responses code",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3319",
              "format": "short",
              "label": "responses",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:3320",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 13
          },
          "hiddenSeries": false,
          "id": 24,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_http_responses_total{proxy=~\"$backend\", code=~\"$code\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ code }}"}} {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - HTTP responses code",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3237",
              "format": "short",
              "label": "responses",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:3238",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 24,
            "x": 0,
            "y": 26
          },
          "height": "400px",
          "hiddenSeries": false,
          "id": 64,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_http_responses_total{proxy=~\"$backend\",server=~\"$server\",code=~\"$code\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ code }}"}} {{"{{ proxy }}"}} {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - HTTP responses code",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:7223",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:7224",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "repeat": null,
      "title": "Responses by HTTP code",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 41
      },
      "id": 157,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 8
          },
          "hiddenSeries": false,
          "id": 45,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:3825",
              "alias": "/.*Denied.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_frontend_sessions_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_frontend_denied_sessions_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Denied {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "haproxy_frontend_current_sessions{proxy=~\"$frontend\",instance=\"$host\"}",
              "hide": true,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current active {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Number of sessions",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3843",
              "format": "short",
              "label": "- denied / + total",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:3844",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 8
          },
          "hiddenSeries": false,
          "id": 30,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_sessions_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_backend_current_sessions{proxy=~\"$backend\",instance=\"$host\"}",
              "hide": true,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current active {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Number of sessions",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3456",
              "format": "short",
              "label": "total",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:3457",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 1,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 21
          },
          "hiddenSeries": false,
          "id": 34,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "maxPerRow": 4,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:4402",
              "alias": "/.*Limit.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_max_sessions{proxy=~\"$backend\",instance=\"$host\"}",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Max {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_backend_limit_sessions{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Limit {{"{{ proxy }}"}}",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Maximum observed number of active sessions and limit",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:4415",
              "format": "short",
              "label": "sessions",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:4416",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 21
          },
          "hiddenSeries": false,
          "id": 51,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 4,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:4155",
              "alias": "/.*Limit.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_frontend_max_sessions{proxy=~\"$frontend\",instance=\"$host\"}",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Max {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_frontend_limit_sessions{proxy=~\"$frontend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Limit {{"{{ proxy }}"}}",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Maximum observed number of active sessions and limit",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:4168",
              "format": "short",
              "label": "sessions",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:4169",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 34
          },
          "hiddenSeries": false,
          "id": 33,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_max_session_rate{proxy=~\"$backend\",instance=\"$host\"}",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Max {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Maximum observed number of sessions per second",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:715",
              "format": "short",
              "label": "sessions",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:716",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 34
          },
          "hiddenSeries": false,
          "id": 69,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 3,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:937",
              "alias": "/.*Limit.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_frontend_max_session_rate{proxy=~\"$frontend\",instance=\"$host\"}",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Max {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_frontend_limit_session_rate{proxy=~\"$frontend\",instance=\"$host\"}",
              "interval": "",
              "intervalFactor": 2,
              "legendFormat": "Limit {{"{{ proxy }}"}}",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Maximum observed number of sessions per second and limit",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:950",
              "format": "short",
              "label": "sessions",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:951",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 47
          },
          "hiddenSeries": false,
          "id": 117,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_failed_header_rewriting_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Failed header rewriting warnings",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1200",
              "format": "short",
              "label": "sessions",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:1201",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 47
          },
          "hiddenSeries": false,
          "id": 119,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_frontend_failed_header_rewriting_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Failed header rewriting warnings",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1374",
              "format": "short",
              "label": "sessions",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:1375",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 60
          },
          "hiddenSeries": false,
          "id": 124,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:4477",
              "alias": "/.*Error.*/",
              "transform": "negative-Y"
            },
            {
              "$$hashKey": "object:4478",
              "alias": "/.*Denied.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_last_session_seconds{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back -  Last session assigned",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:4497",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:4498",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 60
          },
          "hiddenSeries": false,
          "id": 120,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 4,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_frontend_http_requests_rate_max{proxy=~\"$frontend\",instance=\"$host\"}",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Maximum observed number of HTTP requests per second",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1940",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:1941",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 73
          },
          "hiddenSeries": false,
          "id": 128,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:1883",
              "alias": "/.*By server.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_client_aborts_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "By client {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_backend_server_aborts_total{proxy=~\"$frontend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "By server {{"{{ proxy }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back -  Data transfers aborted",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1792",
              "format": "s",
              "label": "- server / + client",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:1793",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 73
          },
          "hiddenSeries": false,
          "id": 146,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "alias": "/.*server.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_client_aborts_total{proxy=~\"$frontend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "By client {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_server_server_aborts_total{proxy=~\"$frontend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "By server {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server -  Data transfers aborted",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "repeat": null,
      "title": "Sessions",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 42
      },
      "id": 158,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 0,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 0,
            "y": 9
          },
          "hiddenSeries": false,
          "id": 38,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_check_up_down_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - UP->DOWN transitions",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 1,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:2029",
              "format": "short",
              "label": "transitions",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:2030",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 12,
            "y": 9
          },
          "hiddenSeries": false,
          "id": 39,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_weight{proxy=~\"$backend\", instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Service weight",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 1,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:2191",
              "format": "short",
              "label": "weight",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:2192",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 12,
            "y": 23
          },
          "hiddenSeries": false,
          "id": 220,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_uweight{proxy=~\"$backend\", instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Service user weight",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 1,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:2191",
              "format": "short",
              "label": "weight",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:2192",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "repeat": null,
      "title": "Health and Weight",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 43
      },
      "id": 159,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 0,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 12,
            "x": 0,
            "y": 10
          },
          "hiddenSeries": false,
          "id": 121,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:2754",
              "alias": "/.*Hits.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_frontend_http_cache_lookups_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Lookups {{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_frontend_http_cache_hits_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Hits {{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Cache lookups / hits",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 1,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:2767",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:2768",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 0,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 12,
            "x": 12,
            "y": 10
          },
          "hiddenSeries": false,
          "id": 139,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:2587",
              "alias": "/.*Hits.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_http_cache_lookups_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Lookups {{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_backend_http_cache_hits_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Hits {{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Cache lookups / hits",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 1,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:2600",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:2601",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 0,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 12,
            "x": 0,
            "y": 21
          },
          "hiddenSeries": false,
          "id": 122,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:3160",
              "alias": "/.*emitted.*/",
              "transform": "negative-Y"
            },
            {
              "$$hashKey": "object:3161",
              "alias": "/.*bypassed.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_frontend_http_comp_bytes_in_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Bytes fed {{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_frontend_http_comp_bytes_out_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Bytes emitted {{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_frontend_http_comp_bytes_bypassed_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Bytes bypassed  {{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Compressor",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 1,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3180",
              "format": "bytes",
              "label": "- emitted / + bypasses + fed",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:3181",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 0,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 12,
            "x": 12,
            "y": 21
          },
          "hiddenSeries": false,
          "id": 140,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:2829",
              "alias": "/.*emitted.*/",
              "transform": "negative-Y"
            },
            {
              "$$hashKey": "object:2830",
              "alias": "/.*bypassed.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_http_comp_bytes_in_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Bytes fed {{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_backend_http_comp_bytes_out_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Bytes emitted {{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_backend_http_comp_bytes_bypassed_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Bytes bypassed  {{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Compressor",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 1,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:2849",
              "format": "bytes",
              "label": "- emitted / + bypasses + fed",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:2850",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 0,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 12,
            "x": 0,
            "y": 32
          },
          "hiddenSeries": false,
          "id": 123,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_frontend_http_comp_responses_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Responses compressed",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 1,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3330",
              "format": "short",
              "label": "responses",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:3331",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 0,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 12,
            "x": 12,
            "y": 32
          },
          "hiddenSeries": false,
          "id": 141,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_http_comp_responses_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} ",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Responses compressed",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 1,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3248",
              "format": "short",
              "label": "responses",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:3249",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "repeat": null,
      "title": "Cache / Compressor",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 44
      },
      "id": 160,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 12,
            "x": 0,
            "y": 11
          },
          "hiddenSeries": false,
          "id": 113,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_frontend_status{proxy=~\"$frontend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Status",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:4278",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:4279",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 12,
            "x": 12,
            "y": 11
          },
          "hiddenSeries": false,
          "id": 112,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_status{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Status",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:4144",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:4145",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 1,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 12,
            "x": 0,
            "y": 22
          },
          "hiddenSeries": false,
          "id": 205,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_frontend_internal_errors_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Front - Internal errors",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1409",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:1410",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 1,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 12,
            "x": 12,
            "y": 22
          },
          "hiddenSeries": false,
          "id": 171,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_backend_internal_errors_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Internal errors",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1409",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:1410",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 12,
            "x": 0,
            "y": 33
          },
          "hiddenSeries": false,
          "id": 173,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_active_servers{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Active servers",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:227",
              "format": "short",
              "label": "- backup / + active",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:228",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 12,
            "x": 12,
            "y": 33
          },
          "hiddenSeries": false,
          "id": 208,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_backend_backup_servers{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "",
              "legendFormat": "{{"{{ proxy }}"}}",
              "refId": "B"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Back - Backup servers",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:227",
              "format": "short",
              "label": "- backup / + active",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:228",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "repeat": null,
      "title": "Status",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 45
      },
      "id": 197,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 1,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 8,
            "w": 24,
            "x": 0,
            "y": 12
          },
          "hiddenSeries": false,
          "id": 129,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:5034",
              "alias": "/.*OUT.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_bytes_in_total{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])*8",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "IN {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_server_bytes_out_total{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])*8",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "OUT {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Incoming / Outgoing bytes",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:5047",
              "format": "bits",
              "label": "- out / + in",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:5048",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 9,
            "w": 24,
            "x": 0,
            "y": 20
          },
          "hiddenSeries": false,
          "id": 219,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:16701",
              "alias": "/.*Estimated.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_used_connections_current{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "In use {{"{{ proxy }}"}} / {{"{{ server }}"}}{{"{{ proxy }}"}}",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_server_idle_connections_current{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Estimated {{"{{ proxy }}"}} / {{"{{ server }}"}}{{"{{ proxy }}"}}",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Connections",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:356",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:357",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 8,
            "w": 24,
            "x": 0,
            "y": 29
          },
          "hiddenSeries": false,
          "id": 130,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:5195",
              "alias": "/.*Error.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_connection_attempts_total{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Attempts {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_server_connection_errors_total{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Error {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Connections attempts / error",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:5208",
              "format": "short",
              "label": "- error / + attempts",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:5209",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 9,
            "w": 24,
            "x": 0,
            "y": 37
          },
          "hiddenSeries": false,
          "id": 179,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_connection_reuses_total{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Reuses {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Number of connections reuses",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:5208",
              "format": "short",
              "label": "reuses",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:5209",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 9,
            "w": 24,
            "x": 0,
            "y": 46
          },
          "hiddenSeries": false,
          "id": 186,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:16701",
              "alias": "/.*limit.*/",
              "fill": 0
            },
            {
              "$$hashKey": "object:1142",
              "alias": "/.*unsafe.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_idle_connections_current{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Available idle connections {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "refId": "C",
              "step": 240
            },
            {
              "expr": "rate(haproxy_server_idle_connections_limit{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Limit available idle connections {{"{{ proxy }}"}} / {{"{{ server }}"}}{{"{{ proxy }}"}}",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_server_safe_idle_connections_current{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Number of safe idle connections {{"{{ proxy }}"}} / {{"{{ server }}"}}{{"{{ proxy }}"}}",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_server_unsafe_idle_connections_current{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Number of unsafe idle connections {{"{{ proxy }}"}} / {{"{{ server }}"}}{{"{{ proxy }}"}}",
              "refId": "D",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Idle connections",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:356",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:357",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "title": "By server - Throughtput  / Connections",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 46
      },
      "id": 201,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 9,
            "w": 24,
            "x": 0,
            "y": 13
          },
          "hiddenSeries": false,
          "id": 187,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "hideEmpty": false,
            "hideZero": false,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_responses_denied_total{proxy=~\"$backend\", instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Denied {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - HTTP responses denied",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1906",
              "format": "short",
              "label": "denied",
              "logBase": 1,
              "max": null,
              "min": "0",
              "show": true
            },
            {
              "$$hashKey": "object:1907",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 7,
            "w": 24,
            "x": 0,
            "y": 22
          },
          "hiddenSeries": false,
          "id": 71,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "hideEmpty": false,
            "hideZero": false,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:8949",
              "alias": "/.*Error.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_redispatch_warnings_total{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Redispatch {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_server_retry_warnings_total{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Retry {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_server_response_errors_total{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 2,
              "legendFormat": "Error {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Error / Redispatch / Retry",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:8962",
              "format": "short",
              "label": "- error / + redispatch + retry",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:8963",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 7,
            "w": 24,
            "x": 0,
            "y": 29
          },
          "hiddenSeries": false,
          "id": 59,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_current_queue{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Queued requests",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:6645",
              "format": "short",
              "label": "requests",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:6646",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 7,
            "w": 24,
            "x": 0,
            "y": 36
          },
          "hiddenSeries": false,
          "id": 180,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:8708",
              "alias": "/.*Limit.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_max_queue{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Max {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "haproxy_server_queue_limit{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "",
              "legendFormat": "Limit {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "refId": "A"
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Max  queued requests and limit",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:6645",
              "format": "short",
              "label": "requests",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:6646",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "title": "By server - Requests / Responses / Queues",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 47
      },
      "id": 193,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 7,
            "w": 24,
            "x": 0,
            "y": 14
          },
          "hiddenSeries": false,
          "id": 135,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_connect_time_average_seconds{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "haproxy_backend_current_queue",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Average connection time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:9202",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:9203",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 7,
            "w": 24,
            "x": 0,
            "y": 21
          },
          "hiddenSeries": false,
          "id": 190,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_max_connect_time_seconds{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "haproxy_backend_current_queue",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Max connection time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:9202",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:9203",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 7,
            "w": 24,
            "x": 0,
            "y": 28
          },
          "hiddenSeries": false,
          "id": 183,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_response_time_average_seconds{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "haproxy_backend_current_queue",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Average response time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:9498",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:9499",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 7,
            "w": 24,
            "x": 0,
            "y": 35
          },
          "hiddenSeries": false,
          "id": 189,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_max_response_time_seconds{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "haproxy_backend_current_queue",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Max response time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:9498",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:9499",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 8,
            "w": 24,
            "x": 0,
            "y": 42
          },
          "hiddenSeries": false,
          "id": 184,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_total_time_average_seconds{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "haproxy_backend_current_queue",
              "refId": "D",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Average total time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:10006",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:10007",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 9,
            "w": 24,
            "x": 0,
            "y": 50
          },
          "hiddenSeries": false,
          "id": 188,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_max_total_time_seconds{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "haproxy_backend_current_queue",
              "refId": "D",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Max total time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:10006",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:10007",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 9,
            "w": 24,
            "x": 0,
            "y": 59
          },
          "hiddenSeries": false,
          "id": 133,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_queue_time_average_seconds{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Average queue time for last 1024 successful connections",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:10428",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:10429",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 0,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 8,
            "w": 24,
            "x": 0,
            "y": 68
          },
          "hiddenSeries": false,
          "id": 191,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_max_queue_time_seconds{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Max queue time",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "cumulative"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:10428",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:10429",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "title": "By server - Times",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 48
      },
      "id": 214,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 8,
            "w": 24,
            "x": 0,
            "y": 15
          },
          "hiddenSeries": false,
          "id": 61,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:7407",
              "alias": "/.*Limit.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_current_sessions{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "hide": true,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_server_sessions_total{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_server_max_sessions{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])",
              "hide": true,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Max {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "C",
              "step": 240
            },
            {
              "expr": "haproxy_server_limit_sessions{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Limit {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "D",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Number of active sessions",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:7420",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:7421",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 9,
            "w": 24,
            "x": 0,
            "y": 23
          },
          "hiddenSeries": false,
          "id": 137,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null as zero",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:7532",
              "alias": "/.*Error.*/",
              "transform": "negative-Y"
            },
            {
              "$$hashKey": "object:7533",
              "alias": "/.*Denied.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_failed_header_rewriting_total{proxy=~\"$frontend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} /  {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Failed header rewriting warnings",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:7552",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:7553",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 8,
            "w": 24,
            "x": 0,
            "y": 32
          },
          "hiddenSeries": false,
          "id": 60,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "dataLinks": []
          },
          "percentage": false,
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:7616",
              "alias": "/.*Limit.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_max_session_rate{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Rate {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_server_limit_session_rate{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Limit {{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Number of sessions per second and limit",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:7629",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:7630",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "title": "By server - Sessions",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 49
      },
      "id": 203,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 8,
            "w": 24,
            "x": 0,
            "y": 16
          },
          "hiddenSeries": false,
          "id": 73,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_weight{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Service weight",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:11421",
              "format": "none",
              "label": "weight",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:11422",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "Server's user weight, or sum of active servers' user weights for a backend",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 8,
            "w": 24,
            "x": 0,
            "y": 24
          },
          "hiddenSeries": false,
          "id": 215,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_uweight{proxy=~\"$backend\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Users weight",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:11421",
              "format": "none",
              "label": "weight",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:11422",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 1,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 8,
            "w": 24,
            "x": 0,
            "y": 32
          },
          "hiddenSeries": false,
          "id": 56,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_check_up_down_total{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - UP->DOWN transitions",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:11777",
              "format": "none",
              "label": "transitions",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:11778",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 1,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 24,
            "x": 0,
            "y": 40
          },
          "hiddenSeries": false,
          "id": 185,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_check_failures_total{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Checks failures",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:12042",
              "format": "none",
              "label": "failures",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:12043",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 1,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 11,
            "w": 24,
            "x": 0,
            "y": 51
          },
          "hiddenSeries": false,
          "id": 204,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_check_duration_seconds{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Checks duration",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:12042",
              "format": "s",
              "label": "",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:12043",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 1,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 10,
            "w": 24,
            "x": 0,
            "y": 62
          },
          "hiddenSeries": false,
          "id": 90,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_current_throttle{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server -Throttle percentage",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:11572",
              "format": "percent",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:11573",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 0,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 7,
            "w": 24,
            "x": 0,
            "y": 72
          },
          "hiddenSeries": false,
          "id": 144,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_last_session_seconds{proxy=~\"$backend\", server=~\"$server\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Number of seconds since last session assigned to server/backend",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 1,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:12374",
              "format": "s",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:12375",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "title": "By server - Health and Weight",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 50
      },
      "id": 207,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 8,
            "w": 24,
            "x": 0,
            "y": 17
          },
          "hiddenSeries": false,
          "id": 145,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:6492",
              "alias": "/.*Error.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_server_status{proxy=~\"$backend\",server=~\"$server\",instance=\"$host\"}",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}} / {{"{{ server }}"}}",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Status",
          "tooltip": {
            "msResolution": true,
            "shared": false,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:6505",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:6506",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 1,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 7,
            "w": 24,
            "x": 0,
            "y": 25
          },
          "hiddenSeries": false,
          "id": 172,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_server_internal_errors_total{proxy=~\"$backend\",instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "{{"{{ proxy }}"}}",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Internal errors",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1409",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:1410",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "description": "Number of failed DNS resolutions in current worker process since started",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 1,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 7,
            "w": 24,
            "x": 0,
            "y": 32
          },
          "hiddenSeries": false,
          "id": 216,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": true,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 2,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_process_failed_resolutions{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Failed",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Server - Failed DNS resolutions",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:1409",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:1410",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "title": "By server - Status",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 51
      },
      "id": 166,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 0,
            "y": 18
          },
          "hiddenSeries": false,
          "id": 101,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:13045",
              "alias": "/.*Configured.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_current_session_rate{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current number of sessions per second over last elapsed second",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_process_limit_session_rate{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Configured maximum number of sessions per second",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "haproxy_process_max_session_rate{instance=\"$host\"}",
              "hide": true,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Maximum observed number of sessions per second",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Sessions over last second",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:13058",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:13059",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 12,
            "y": 18
          },
          "hiddenSeries": false,
          "id": 96,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_current_connections{instance=\"$host\"}",
              "hide": true,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Number of active sessions",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_process_connections_total{instance=\"$host\"}[$__rate_interval])",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total number of created sessions",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_process_requests_total{instance=\"$host\"}[$__rate_interval])",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total number of requests (TCP or HTTP)",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Total sessions / requests",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:13140",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:13141",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 0,
            "y": 32
          },
          "hiddenSeries": false,
          "id": 100,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:13202",
              "alias": "/.*Configured.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_current_connection_rate{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current number of connections per second over last elapsed second",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_process_limit_connection_rate{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Configured maximum number of connections per second.",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "haproxy_process_max_connection_rate{instance=\"$host\"}",
              "hide": true,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Maximum observed number of connections per second",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Connections over last second",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:13215",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:13216",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 14,
            "w": 12,
            "x": 12,
            "y": 32
          },
          "hiddenSeries": false,
          "id": 95,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:13277",
              "alias": "/.*Initial.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_max_connections{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Maximum number of concurrent connections",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_process_hard_max_connections{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Initial Maximum number of concurrent connections",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Max connections",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:13290",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:13291",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "repeat": null,
      "title": "Process Connections / Sessions / Requests",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 52
      },
      "id": 167,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 10,
            "w": 12,
            "x": 0,
            "y": 19
          },
          "hiddenSeries": false,
          "id": 106,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:13482",
              "alias": "/.*Configured.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_current_zlib_memory{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current memory used for zlib in bytes",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_process_max_zlib_memory{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Configured maximum amount of memory for zlib in bytes",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Compression memory",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:13495",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:13496",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 10,
            "w": 12,
            "x": 12,
            "y": 19
          },
          "hiddenSeries": false,
          "id": 105,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:13571",
              "alias": "/.*Configured.*/",
              "fill": 0,
              "stack": false
            },
            {
              "$$hashKey": "object:13856",
              "alias": "/.*before.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_process_http_comp_bytes_in_total{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Number of bytes per second over last elapsed second, before http compression",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_process_http_comp_bytes_out_total{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Number of bytes per second over last elapsed second, after http compression",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "haproxy_process_limit_http_comp{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Configured maximum input compression rate in bytes",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Compression",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:13591",
              "format": "short",
              "label": "- before / + after",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:13592",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "repeat": null,
      "title": "Process Compression",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 53
      },
      "id": 168,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 12,
            "w": 12,
            "x": 0,
            "y": 20
          },
          "hiddenSeries": false,
          "id": 104,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:13920",
              "alias": "/.*misses.*/",
              "transform": "negative-Y"
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_process_ssl_cache_lookups_total{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total number of SSL session cache lookups",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_process_ssl_cache_misses_total{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total number of SSL session cache misses",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - SSL cache",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:13933",
              "format": "short",
              "label": "- misses / + lookups",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:13934",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 12,
            "w": 12,
            "x": 12,
            "y": 20
          },
          "hiddenSeries": false,
          "id": 103,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:13995",
              "alias": "/.*backend.*/",
              "transform": "negative-Y"
            },
            {
              "$$hashKey": "object:13996",
              "alias": "/.*Maximum.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_current_frontend_ssl_key_rate{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current frontend SSL Key computation per second over last elapsed second",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_process_max_frontend_ssl_key_rate{instance=\"$host\"}",
              "hide": true,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Maximum observed frontend SSL Key computation per second",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "haproxy_process_current_backend_ssl_key_rate{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current backend SSL Key computation per second over last elapsed second",
              "metric": "",
              "refId": "D",
              "step": 240
            },
            {
              "expr": "haproxy_process_max_backend_ssl_key_rate{instance=\"$host\"}",
              "hide": true,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Maximum observed backend SSL Key computation per second",
              "metric": "",
              "refId": "E",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - SSL key rate",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:14015",
              "format": "short",
              "label": "- backend / + frontend",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": true
            },
            {
              "$$hashKey": "object:14016",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 12,
            "w": 12,
            "x": 0,
            "y": 32
          },
          "hiddenSeries": false,
          "id": 102,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:14311",
              "alias": "/.*Maximum.*/",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_current_ssl_rate{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current number of SSL sessions per second over last elapsed second",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_process_limit_ssl_rate{instance=\"$host\"}",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Configured maximum number of SSL sessions per second",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "haproxy_process_max_ssl_rate{instance=\"$host\"}",
              "hide": true,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Maximum observed number of SSL sessions per second",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - SSL rate",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:14324",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:14325",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 12,
            "w": 12,
            "x": 12,
            "y": 32
          },
          "hiddenSeries": false,
          "id": 98,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:14386",
              "alias": "/.*Maximum*./",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_current_ssl_connections{instance=\"$host\"}",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current number of opened SSL connections",
              "metric": "",
              "refId": "D",
              "step": 240
            },
            {
              "expr": "haproxy_process_max_ssl_connections{instance=\"$host\"}",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Configured maximum number of concurrent SSL connections",
              "metric": "",
              "refId": "E",
              "step": 240
            },
            {
              "expr": "rate(haproxy_process_ssl_connections_total{instance=\"$host\"}[$__rate_interval])",
              "hide": true,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total number of opened SSL connections",
              "metric": "",
              "refId": "F",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - SSL connections",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:14399",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:14400",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 12,
            "w": 12,
            "x": 0,
            "y": 44
          },
          "hiddenSeries": false,
          "id": 150,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_frontend_ssl_reuse{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "SSL session reuse ratio (percent)",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - SSL reuse ratio",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:14474",
              "decimals": null,
              "format": "percent",
              "label": "ratio",
              "logBase": 1,
              "max": null,
              "min": "0",
              "show": true
            },
            {
              "$$hashKey": "object:14475",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "repeat": null,
      "title": "Process SSL",
      "type": "row"
    },
    {
      "collapsed": true,
      "datasource": "{{.DatasourceName}}",
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 54
      },
      "id": 169,
      "panels": [
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 55
          },
          "hiddenSeries": false,
          "id": 87,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "alias": "/.*limit+./",
              "fill": 0
            }
          ],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_max_memory_bytes{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Per-process memory limit (in bytes); 0=unset",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_process_pool_allocated_bytes{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total amount of memory allocated in pools (in bytes)",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_process_pool_used_bytes{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total amount of memory used in pools (in bytes)",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Memory",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:980",
              "format": "bytes",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:981",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 55
          },
          "hiddenSeries": false,
          "id": 107,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_process_current_tasks{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current number of tasks",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_process_current_run_queue{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current number of tasks in the run-queue",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "haproxy_process_stopping{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Non zero means stopping in progress",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Tasks",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:14780",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:14781",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 68
          },
          "hiddenSeries": false,
          "id": 89,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_max_fds{instance=\"$host\"}",
              "hide": false,
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Maximum number of open file descriptors; 0=unset",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_process_max_sockets{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Maximum numer of open sockets",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Maximum open files / sockets",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:14848",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:14849",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 68
          },
          "hiddenSeries": false,
          "id": 99,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:14910",
              "alias": "/.*Configured.*/",
              "fill": 0,
              "stack": false
            }
          ],
          "spaceLength": 10,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_max_pipes{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Configured maximum number of pipes",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_process_pipes_used_total{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Number of pipes in used",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_process_pipes_free_total{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Number of pipes unused",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Pipes",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:14930",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:14931",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 81
          },
          "hiddenSeries": false,
          "id": 86,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_nbthread{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Threads",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_process_nbproc{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Processes",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Configured threads / processes",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:15002",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:15003",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 81
          },
          "hiddenSeries": false,
          "id": 88,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_pool_failures_total{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total number of failed pool allocations",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Pool allocations",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:15070",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:15071",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 94
          },
          "hiddenSeries": false,
          "id": 108,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_idle_time_percent{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Idle to total ratio over last sample (percent)",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Idle",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "format": "percent",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 94
          },
          "hiddenSeries": false,
          "id": 109,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_process_jobs{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current number of active jobs (listeners, sessions, open devices)",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_process_unstoppable_jobs{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current number of active jobs that can't be stopped during a soft stop",
              "metric": "",
              "refId": "B",
              "step": 240
            },
            {
              "expr": "rate(haproxy_process_listeners{instance=\"$host\"}[$__rate_interval])",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current number of active listeners",
              "metric": "",
              "refId": "C",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Jobs",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:15138",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:15139",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 0,
            "y": 107
          },
          "hiddenSeries": false,
          "id": 110,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_active_peers{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current number of active peers",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_process_connected_peers{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Current number of connected peers",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Peers",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:15206",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:15207",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 107
          },
          "hiddenSeries": false,
          "id": 111,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_dropped_logs_total{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total number of dropped logs",
              "metric": "",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "haproxy_process_recv_logs_total{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Total number of log messages received by log-forwarding listeners",
              "metric": "",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Logs",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:15274",
              "format": "short",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:15275",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "decimals": 1,
          "description": "",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 15,
            "w": 12,
            "x": 0,
            "y": 120
          },
          "hiddenSeries": false,
          "id": 217,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 2,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [],
          "spaceLength": 10,
          "stack": false,
          "steppedLine": false,
          "targets": [
            {
              "expr": "rate(haproxy_process_bytes_out_total{instance=\"$host\"}[$__rate_interval])*8",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Emitted by current worker {{"{{ proxy }}"}}",
              "metric": "haproxy_backend_",
              "refId": "A",
              "step": 240
            },
            {
              "expr": "rate(haproxy_process_spliced_bytes_out_total{instance=\"$host\"}[$__rate_interval])*8",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "Emitted by current worker through a kernel pipe {{"{{ proxy }}"}}",
              "metric": "haproxy_backend_",
              "refId": "B",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Bytes out",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 0,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:3132",
              "format": "bits",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": null,
              "show": false
            },
            {
              "$$hashKey": "object:3133",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        },
        {
          "aliasColors": {},
          "bars": false,
          "dashLength": 10,
          "dashes": false,
          "datasource": "{{.DatasourceName}}",
          "editable": true,
          "error": false,
          "fieldConfig": {
            "defaults": {
              "custom": {},
              "links": []
            },
            "overrides": []
          },
          "fill": 2,
          "fillGradient": 0,
          "grid": {},
          "gridPos": {
            "h": 13,
            "w": 12,
            "x": 12,
            "y": 120
          },
          "hiddenSeries": false,
          "id": 218,
          "legend": {
            "alignAsTable": true,
            "avg": true,
            "current": true,
            "max": true,
            "min": true,
            "rightSide": false,
            "show": true,
            "total": false,
            "values": true
          },
          "lines": true,
          "linewidth": 1,
          "links": [],
          "maxPerRow": 1,
          "nullPointMode": "null",
          "options": {
            "alertThreshold": true
          },
          "percentage": false,
          "pluginVersion": "7.3.7",
          "pointradius": 5,
          "points": false,
          "renderer": "flot",
          "seriesOverrides": [
            {
              "$$hashKey": "object:14910",
              "alias": "/.*Configured.*/",
              "fill": 0,
              "stack": false
            }
          ],
          "spaceLength": 10,
          "stack": true,
          "steppedLine": false,
          "targets": [
            {
              "expr": "haproxy_process_uptime_seconds{instance=\"$host\"}",
              "interval": "$interval",
              "intervalFactor": 1,
              "legendFormat": "How long ago this worker process was started",
              "metric": "",
              "refId": "A",
              "step": 240
            }
          ],
          "thresholds": [],
          "timeFrom": null,
          "timeRegions": [],
          "timeShift": null,
          "title": "Process - Uptime",
          "tooltip": {
            "msResolution": true,
            "shared": true,
            "sort": 2,
            "value_type": "individual"
          },
          "type": "graph",
          "xaxis": {
            "buckets": null,
            "mode": "time",
            "name": null,
            "show": true,
            "values": []
          },
          "yaxes": [
            {
              "$$hashKey": "object:14930",
              "format": "s",
              "label": "counter",
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": true
            },
            {
              "$$hashKey": "object:14931",
              "format": "short",
              "label": null,
              "logBase": 1,
              "max": null,
              "min": 0,
              "show": false
            }
          ],
          "yaxis": {
            "align": false,
            "alignLevel": null
          }
        }
      ],
      "repeat": null,
      "title": "Process Misc",
      "type": "row"
    }
  ],
  "refresh": "5m",
  "schemaVersion": 26,
  "style": "dark",
  "tags": [
    "haproxy",
    "servers"
  ],
  "templating": {
    "list": [
      {
        "allValue": null,
        "current": {},
        "datasource": "{{.DatasourceName}}",
        "definition": "",
        "error": null,
        "hide": 0,
        "includeAll": false,
        "label": "Host",
        "multi": false,
        "name": "host",
        "options": [],
        "query": "label_values(haproxy_process_nbproc,instance)",
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 1,
        "tagValuesQuery": null,
        "tags": [],
        "tagsQuery": null,
        "type": "query",
        "useTags": false
      },
      {
        "allValue": null,
        "current": {},
        "datasource": "{{.DatasourceName}}",
        "definition": "",
        "error": null,
        "hide": 0,
        "includeAll": true,
        "label": "Backend",
        "multi": true,
        "name": "backend",
        "options": [],
        "query": "label_values(haproxy_backend_status{instance=\"$host\"}, proxy)",
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 1,
        "tagValuesQuery": null,
        "tags": [],
        "tagsQuery": null,
        "type": "query",
        "useTags": false
      },
      {
        "allValue": null,
        "current": {},
        "datasource": "{{.DatasourceName}}",
        "definition": "",
        "error": null,
        "hide": 0,
        "includeAll": true,
        "label": "Frontend",
        "multi": true,
        "name": "frontend",
        "options": [],
        "query": "label_values(haproxy_frontend_status{instance=\"$host\"}, proxy)",
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 1,
        "tagValuesQuery": null,
        "tags": [],
        "tagsQuery": null,
        "type": "query",
        "useTags": true
      },
      {
        "allValue": null,
        "current": {},
        "datasource": "{{.DatasourceName}}",
        "definition": "",
        "error": null,
        "hide": 0,
        "includeAll": true,
        "label": "Server",
        "multi": true,
        "name": "server",
        "options": [],
        "query": "label_values(haproxy_server_status{instance=\"$host\"}, server)",
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 1,
        "tagValuesQuery": null,
        "tags": [],
        "tagsQuery": null,
        "type": "query",
        "useTags": false
      },
      {
        "allValue": null,
        "current": {},
        "datasource": "{{.DatasourceName}}",
        "definition": "",
        "error": null,
        "hide": 0,
        "includeAll": true,
        "label": "HTTP Code",
        "multi": true,
        "name": "code",
        "options": [],
        "query": "label_values(haproxy_server_http_responses_total{instance=\"$host\"}, code)",
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 1,
        "tagValuesQuery": null,
        "tags": [],
        "tagsQuery": null,
        "type": "query",
        "useTags": false
      },
      {
        "auto": true,
        "auto_count": 30,
        "auto_min": "10s",
        "current": {
          "selected": false,
          "text": "30s",
          "value": "30s"
        },
        "error": null,
        "hide": 0,
        "label": "Interval",
        "name": "interval",
        "options": [
          {
            "selected": false,
            "text": "auto",
            "value": "$__auto_interval_interval"
          },
          {
            "selected": true,
            "text": "30s",
            "value": "30s"
          },
          {
            "selected": false,
            "text": "1m",
            "value": "1m"
          },
          {
            "selected": false,
            "text": "5m",
            "value": "5m"
          },
          {
            "selected": false,
            "text": "1h",
            "value": "1h"
          },
          {
            "selected": false,
            "text": "6h",
            "value": "6h"
          },
          {
            "selected": false,
            "text": "1d",
            "value": "1d"
          }
        ],
        "query": "30s,1m,5m,1h,6h,1d",
        "refresh": 2,
        "skipUrlSync": false,
        "type": "interval"
      }
    ]
  },
  "time": {
    "from": "now-24h",
    "to": "now"
  },
  "timepicker": {
    "refresh_intervals": [
      "5s",
      "10s",
      "30s",
      "1m",
      "5m",
      "15m",
      "30m",
      "1h",
      "2h",
      "1d"
    ],
    "time_options": [
      "5m",
      "15m",
      "1h",
      "6h",
      "12h",
      "24h",
      "2d",
      "7d",
      "30d"
    ]
  },
  "timezone": "browser",
  "title": "HAProxy 2 {{.Title}}",
  "uid": "rEqu1u5ue",
  "version": 9
}`
