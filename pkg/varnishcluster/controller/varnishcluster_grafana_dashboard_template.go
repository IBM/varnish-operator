package controller

var grafanaDashboardTemplate = `
{
  "annotations": {
    "list": [
      {
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
  "editable": true,
  "gnetId": null,
  "graphTooltip": 1,
  "iteration": 1563214363724,
  "links": [],
  "panels": [
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 0
      },
      "id": 19,
      "panels": [],
      "repeat": "cluster",
      "title": "Real time metrics",
      "type": "row"
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
      "decimals": null,
      "description": "",
      "editable": true,
      "error": false,
      "format": "percentunit",
      "gauge": {
        "maxValue": 1,
        "minValue": 0,
        "show": true,
        "thresholdLabels": false,
        "thresholdMarkers": false
      },
      "gridPos": {
        "h": 4,
        "w": 2,
        "x": 0,
        "y": 1
      },
      "height": "",
      "hideTimeOverride": false,
      "id": 6,
      "interval": null,
      "links": [],
      "mappingType": 1,
      "mappingTypes": [
        {
          "name": "value to text",
          "value": 1
        },
        {
          "name": "range to text",
          "value": 2
        }
      ],
      "maxDataPoints": 100,
      "nullPointMode": "connected",
      "nullText": null,
      "options": {},
      "postfix": "",
      "postfixFontSize": "50%",
      "prefix": "",
      "prefixFontSize": "100%",
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
      "tableColumn": "Value",
      "targets": [
        {
          "expr": "avg(\n  (\n    rate(varnish_main_cache_hit{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]) / (rate(varnish_main_cache_hit{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]) + rate(varnish_main_cache_miss{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]))\n  )\n)",
          "format": "time_series",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "",
          "metric": "",
          "refId": "A",
          "step": 1800
        }
      ],
      "thresholds": "0,0,0.1",
      "title": "Cache Hit Rate",
      "transparent": true,
      "type": "singlestat",
      "valueFontSize": "70%",
      "valueMaps": [
        {
          "op": "=",
          "text": "N/A",
          "value": "null"
        }
      ],
      "valueName": "current"
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
      "decimals": 1,
      "editable": true,
      "error": false,
      "format": "none",
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
        "x": 2,
        "y": 1
      },
      "id": 9,
      "interval": null,
      "links": [],
      "mappingType": 1,
      "mappingTypes": [
        {
          "name": "value to text",
          "value": 1
        },
        {
          "name": "range to text",
          "value": 2
        }
      ],
      "maxDataPoints": 100,
      "nullPointMode": "connected",
      "nullText": null,
      "options": {},
      "postfix": " / sec",
      "postfixFontSize": "50%",
      "prefix": "",
      "prefixFontSize": "100%",
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
          "expr": "sum(irate(varnish_main_cache_hit{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m])) + sum(irate(varnish_main_cache_miss{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]))",
          "format": "time_series",
          "hide": false,
          "interval": "15s",
          "intervalFactor": 1,
          "legendFormat": "",
          "metric": "",
          "refId": "A",
          "step": 900
        }
      ],
      "thresholds": "",
      "timeFrom": null,
      "title": "Frontend Requests",
      "transparent": true,
      "type": "singlestat",
      "valueFontSize": "70%",
      "valueMaps": [
        {
          "op": "=",
          "text": "N/A",
          "value": "null"
        }
      ],
      "valueName": "current"
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
      "decimals": 1,
      "editable": true,
      "error": false,
      "format": "short",
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
        "x": 4,
        "y": 1
      },
      "id": 11,
      "interval": null,
      "links": [],
      "mappingType": 1,
      "mappingTypes": [
        {
          "name": "value to text",
          "value": 1
        },
        {
          "name": "range to text",
          "value": 2
        }
      ],
      "maxDataPoints": 100,
      "nullPointMode": "connected",
      "nullText": null,
      "options": {},
      "postfix": " / sec",
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
      "tableColumn": "frontend_for",
      "targets": [
        {
          "expr": "sum(irate(varnish_main_sessions_total{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]))",
          "format": "time_series",
          "hide": false,
          "interval": "",
          "intervalFactor": 1,
          "legendFormat": "",
          "refId": "A",
          "step": 900
        }
      ],
      "thresholds": "",
      "title": "Sessions",
      "transparent": true,
      "type": "singlestat",
      "valueFontSize": "70%",
      "valueMaps": [
        {
          "op": "=",
          "text": "N/A",
          "value": "null"
        }
      ],
      "valueName": "current"
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
      "decimals": 1,
      "editable": true,
      "error": false,
      "format": "none",
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
        "x": 6,
        "y": 1
      },
      "height": "",
      "id": 4,
      "interval": null,
      "links": [],
      "mappingType": 1,
      "mappingTypes": [
        {
          "name": "value to text",
          "value": 1
        },
        {
          "name": "range to text",
          "value": 2
        }
      ],
      "maxDataPoints": 100,
      "nullPointMode": "connected",
      "nullText": null,
      "options": {},
      "postfix": " / sec",
      "postfixFontSize": "50%",
      "prefix": "",
      "prefixFontSize": "100%",
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
          "expr": "sum(irate(varnish_main_cache_hit{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]))",
          "format": "time_series",
          "interval": "15s",
          "intervalFactor": 1,
          "legendFormat": "",
          "metric": "varnish_main_backend_req",
          "refId": "A",
          "step": 900
        }
      ],
      "thresholds": "",
      "title": "Cache Hit",
      "transparent": true,
      "type": "singlestat",
      "valueFontSize": "70%",
      "valueMaps": [
        {
          "op": "=",
          "text": "N/A",
          "value": "null"
        }
      ],
      "valueName": "current"
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
      "format": "none",
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
        "x": 8,
        "y": 1
      },
      "id": 16,
      "interval": null,
      "links": [],
      "mappingType": 1,
      "mappingTypes": [
        {
          "name": "value to text",
          "value": 1
        },
        {
          "name": "range to text",
          "value": 2
        }
      ],
      "maxDataPoints": 100,
      "nullPointMode": "connected",
      "nullText": null,
      "options": {},
      "postfix": "",
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
          "expr": "max(varnish_main_vmods{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"})",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A",
          "step": 1800
        }
      ],
      "thresholds": "",
      "title": "Loaded VMODs",
      "transparent": true,
      "type": "singlestat",
      "valueFontSize": "70%",
      "valueMaps": [
        {
          "op": "=",
          "text": "N/A",
          "value": "null"
        }
      ],
      "valueName": "avg"
    },
    {
      "cacheTimeout": null,
      "colorBackground": false,
      "colorValue": true,
      "colors": [
        "rgba(245, 54, 54, 0.9)",
        "rgb(96, 96, 96)",
        "rgba(50, 172, 45, 0.97)"
      ],
      "datasource": "{{.DatasourceName}}",
      "decimals": null,
      "editable": true,
      "error": false,
      "format": "dtdurations",
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
        "x": 10,
        "y": 1
      },
      "id": 18,
      "interval": null,
      "links": [],
      "mappingType": 1,
      "mappingTypes": [
        {
          "name": "value to text",
          "value": 1
        },
        {
          "name": "range to text",
          "value": 2
        }
      ],
      "maxDataPoints": 100,
      "nullPointMode": "connected",
      "nullText": null,
      "options": {},
      "postfix": "",
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
      "tableColumn": "__name__",
      "targets": [
        {
          "expr": "min(varnish_main_uptime{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"})",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A",
          "step": 1800
        }
      ],
      "thresholds": "",
      "title": "Uptime",
      "transparent": true,
      "type": "singlestat",
      "valueFontSize": "50%",
      "valueMaps": [
        {
          "op": "=",
          "text": "N/A",
          "value": "null"
        }
      ],
      "valueName": "current"
    },
    {
      "cacheTimeout": null,
      "colorBackground": false,
      "colorValue": true,
      "colors": [
        "rgba(245, 54, 54, 0.9)",
        "rgb(96, 96, 96)",
        "rgba(50, 172, 45, 0.97)"
      ],
      "datasource": "{{.DatasourceName}}",
      "decimals": null,
      "editable": true,
      "error": false,
      "format": "dtdurations",
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
        "x": 12,
        "y": 1
      },
      "id": 10,
      "interval": null,
      "links": [],
      "mappingType": 1,
      "mappingTypes": [
        {
          "name": "value to text",
          "value": 1
        },
        {
          "name": "range to text",
          "value": 2
        }
      ],
      "maxDataPoints": 100,
      "nullPointMode": "connected",
      "nullText": null,
      "options": {},
      "postfix": "",
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
      "tableColumn": "version",
      "targets": [
        {
          "expr": "varnish_version{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}",
          "format": "table",
          "intervalFactor": 2,
          "legendFormat": "",
          "refId": "A",
          "step": 1800
        }
      ],
      "thresholds": "",
      "title": "Version",
      "transparent": true,
      "type": "singlestat",
      "valueFontSize": "50%",
      "valueMaps": [
        {
          "op": "=",
          "text": "N/A",
          "value": "null"
        }
      ],
      "valueName": "current"
    },
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 5
      },
      "id": 91,
      "panels": [],
      "title": "Varnish requests",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "{{.DatasourceName}}",
      "fill": 1,
      "gridPos": {
        "h": 5,
        "w": 8,
        "x": 0,
        "y": 6
      },
      "id": 27,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "sum(irate(varnish_main_cache_hit{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m])) + sum(irate(varnish_main_cache_miss{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m]))",
          "format": "time_series",
          "hide": false,
          "intervalFactor": 1,
          "legendFormat": "requests",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Varnish Overall Req/sec",
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
          "decimals": 0,
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
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
      "fill": 5,
      "gridPos": {
        "h": 5,
        "w": 8,
        "x": 8,
        "y": 6
      },
      "id": 29,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 0,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": true,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(irate(varnish_main_cache_hit{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m]))",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "cache hits",
          "refId": "A"
        },
        {
          "expr": "sum(irate(varnish_main_cache_miss{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m]))",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "cache misses",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Cache Hit/Miss Reqs/sec",
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
          "decimals": 0,
          "format": "reqps",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
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
      "fill": 1,
      "gridPos": {
        "h": 5,
        "w": 8,
        "x": 16,
        "y": 6
      },
      "id": 41,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "sum(irate(varnish_main_thread_queue_len{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m]))",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "queued requests",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Queued Requests/sec",
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
          "decimals": 0,
          "format": "none",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
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
      "fill": 6,
      "gridPos": {
        "h": 5,
        "w": 8,
        "x": 0,
        "y": 11
      },
      "id": 14,
      "legend": {
        "alignAsTable": true,
        "avg": true,
        "current": false,
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
      "linewidth": 0,
      "links": [],
      "nullPointMode": "connected",
      "options": {},
      "paceLength": 10,
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "alias": "/.*miss/",
          "color": "#890F02"
        },
        {
          "alias": "/.*hit/",
          "color": "#508642"
        }
      ],
      "spaceLength": 10,
      "stack": true,
      "steppedLine": true,
      "targets": [
        {
          "expr": "avg(rate(varnish_main_cache_hit{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]) / (rate(varnish_main_cache_miss{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m])+rate(varnish_main_cache_hit{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m])))",
          "format": "time_series",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "cache hit %",
          "refId": "A",
          "step": 240
        },
        {
          "expr": "avg(1 - (rate(varnish_main_cache_hit{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]) / (rate(varnish_main_cache_miss{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m])+rate(varnish_main_cache_hit{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]))))",
          "format": "time_series",
          "hide": false,
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "cache miss %",
          "refId": "B",
          "step": 240
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Cache hit rate",
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
          "format": "percentunit",
          "label": "",
          "logBase": 1,
          "max": "1",
          "min": "0",
          "show": true
        },
        {
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
      "fill": 1,
      "gridPos": {
        "h": 5,
        "w": 8,
        "x": 8,
        "y": 11
      },
      "id": 31,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "sum(irate(varnish_main_cache_hitpass{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m]))",
          "format": "time_series",
          "hide": true,
          "intervalFactor": 1,
          "legendFormat": "passes",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Cache Pass Reqs/sec",
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
          "decimals": 0,
          "format": "none",
          "label": "",
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
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
      "fill": 1,
      "gridPos": {
        "h": 5,
        "w": 8,
        "x": 16,
        "y": 11
      },
      "id": 39,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "sum(irate(varnish_main_req_dropped{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m]))",
          "format": "time_series",
          "instant": false,
          "intervalFactor": 1,
          "legendFormat": "dropped",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Dropped Requests/sec",
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
          "decimals": 0,
          "format": "none",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 16
      },
      "id": 21,
      "panels": [],
      "repeat": null,
      "title": "Varnish memory",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "{{.DatasourceName}}",
      "fill": 1,
      "gridPos": {
        "h": 5,
        "w": 8,
        "x": 0,
        "y": 17
      },
      "id": 43,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "varnish_sma_g_bytes{service=\"{{.ServiceName}}\", type=\"s0\", namespace=\"{{.Namespace}}\"}",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "{{"{{pod}}"}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Memory",
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
          "format": "bytes",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
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
      "decimals": 2,
      "editable": true,
      "error": false,
      "fill": 4,
      "grid": {},
      "gridPos": {
        "h": 10,
        "w": 8,
        "x": 8,
        "y": 17
      },
      "hideTimeOverride": false,
      "id": 13,
      "legend": {
        "alignAsTable": true,
        "avg": false,
        "current": true,
        "hideZero": false,
        "max": true,
        "min": true,
        "rightSide": false,
        "show": true,
        "sort": null,
        "sortDesc": null,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "connected",
      "options": {},
      "paceLength": 10,
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "alias": "/.*expired/",
          "color": "#82B5D8"
        },
        {
          "alias": "/.*LRU moved/",
          "color": "#F9E2D2"
        },
        {
          "alias": "/.*LRU nuked/",
          "color": "#BF1B00"
        },
        {
          "alias": "/.*mem /",
          "color": "#BA43A9",
          "fill": 0,
          "yaxis": 2
        }
      ],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(irate(varnish_main_n_expired{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]))",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "expired",
          "metric": "varnish_main_n_object",
          "refId": "A",
          "step": 240
        },
        {
          "expr": "sum(irate(varnish_main_n_lru_moved{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]))",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "LRU moved",
          "refId": "C",
          "step": 240
        },
        {
          "expr": "sum(irate(varnish_main_n_lru_nuked{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]))",
          "format": "time_series",
          "interval": "",
          "intervalFactor": 2,
          "legendFormat": "LRU nuked",
          "refId": "B",
          "step": 240
        },
        {
          "expr": "sum(varnish_sma_g_bytes{service=\"{{.ServiceName}}\", type=\"s0\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"})",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "mem used",
          "metric": "",
          "refId": "E",
          "step": 240
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Cache",
      "tooltip": {
        "shared": true,
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
          "format": "none",
          "logBase": 1,
          "max": null,
          "min": 1,
          "show": true
        },
        {
          "format": "bytes",
          "logBase": 1,
          "max": null,
          "min": 1,
          "show": true
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
      "fill": 1,
      "gridPos": {
        "h": 5,
        "w": 8,
        "x": 16,
        "y": 17
      },
      "id": 35,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "sum(irate(varnish_main_n_expired{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m]))",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "expired objects",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Expired Objects/sec",
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
          "decimals": 0,
          "format": "none",
          "label": "",
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
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
      "decimals": null,
      "editable": true,
      "error": false,
      "fill": 5,
      "grid": {},
      "gridPos": {
        "h": 5,
        "w": 8,
        "x": 0,
        "y": 22
      },
      "hideTimeOverride": false,
      "id": 12,
      "legend": {
        "alignAsTable": true,
        "avg": true,
        "current": false,
        "hideEmpty": true,
        "hideZero": true,
        "max": true,
        "min": false,
        "rightSide": false,
        "show": true,
        "sideWidth": 400,
        "sort": null,
        "sortDesc": null,
        "total": true,
        "values": true
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "connected",
      "options": {},
      "paceLength": 10,
      "percentage": false,
      "pointradius": 5,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [
        {
          "alias": "/backend .*/",
          "fill": 7
        },
        {
          "alias": "frontend",
          "color": "#508642",
          "stack": false,
          "zindex": -2
        },
        {
          "alias": "/^backend$/",
          "color": "#CCA300",
          "fill": 0,
          "stack": false,
          "zindex": 1
        }
      ],
      "spaceLength": 10,
      "stack": true,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(irate(varnish_main_s_resp_hdrbytes{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]) + irate(varnish_main_s_resp_bodybytes{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]))",
          "format": "time_series",
          "hide": false,
          "interval": "",
          "intervalFactor": 1,
          "legendFormat": "frontend",
          "refId": "A",
          "step": 60
        },
        {
          "expr": "sum(irate(varnish_backend_beresp_hdrbytes{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]) + irate(varnish_backend_beresp_bodybytes{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]))",
          "format": "time_series",
          "hide": false,
          "interval": "",
          "intervalFactor": 1,
          "legendFormat": "backend",
          "metric": "",
          "refId": "B",
          "step": 60
        },
        {
          "expr": "sum(irate(varnish_backend_beresp_hdrbytes{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m]) + irate(varnish_backend_beresp_bodybytes{service=\"{{.ServiceName}}\", pod=~\"^($varnish_pod).*\", namespace=\"{{.Namespace}}\"}[1m])) by (backend)",
          "format": "time_series",
          "intervalFactor": 2,
          "legendFormat": "backend {{"{{ backend }}"}}",
          "refId": "C",
          "step": 120
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Network",
      "tooltip": {
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
          "format": "bytes",
          "label": "",
          "logBase": 1,
          "max": null,
          "min": "1024",
          "show": true
        },
        {
          "format": "short",
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
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
      "fill": 1,
      "gridPos": {
        "h": 5,
        "w": 8,
        "x": 16,
        "y": 22
      },
      "id": 37,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "sum(irate(varnish_main_n_lru_nuked{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m]))",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "nuked",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Nuked Objects/sec",
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
          "decimals": 0,
          "format": "none",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 27
      },
      "id": 93,
      "panels": [],
      "title": "Varnish internals metrics",
      "type": "row"
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "{{.DatasourceName}}",
      "description": "Number of requests sent to sleep without a worker thread because\n        they found a busy object.",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 8,
        "x": 0,
        "y": 28
      },
      "id": 71,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "sum(irate(varnish_main_busy_sleep{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m]))",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "requests",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Requests sent to sleep sleep per seconds",
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
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
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
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 8,
        "x": 8,
        "y": 28
      },
      "id": 67,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "sum(irate(varnish_main_backend_retry{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m]))",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "retries ",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Backend retries per second",
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
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
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
      "description": "Number of times more threads were needed, but limit was reached in\n        a thread pool. See also parameter thread_pool_max.",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 8,
        "x": 16,
        "y": 28
      },
      "id": 69,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "sum(irate(varnish_main_threads_limited{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m]))",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "threads",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Threads limited per second",
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
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
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
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 8,
        "x": 0,
        "y": 34
      },
      "id": 77,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "sum(irate(varnish_main_sessions{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" }[1m])) by (type)",
          "format": "time_series",
          "hide": false,
          "intervalFactor": 1,
          "legendFormat": "{{"{{type}}"}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Number of sessions",
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
          "min": null,
          "show": true
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
      "description": "Number of times session was queued waiting for a thread. See also\n        parameter thread_queue_limit.",
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 8,
        "x": 8,
        "y": 34
      },
      "id": 73,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "sum(varnish_main_thread_queue_len{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\" })",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "queue",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Length of session queue",
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
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": "0",
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
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
      "fill": 1,
      "gridPos": {
        "h": 6,
        "w": 8,
        "x": 16,
        "y": 34
      },
      "id": 75,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "links": [],
      "nullPointMode": "null",
      "options": {},
      "paceLength": 10,
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
          "expr": "sum(irate(varnish_main_fetch{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\", namespace=\"{{.Namespace}}\" }[1m])) by (type)",
          "format": "time_series",
          "intervalFactor": 1,
          "legendFormat": "{{"{{type}}"}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Number of fetches per second",
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
          "min": null,
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    }
  ],
  "refresh": "5s",
  "schemaVersion": 18,
  "style": "dark",
  "tags": [],
  "templating": {
    "list": [
      {
        "allValue": null,
        "current": {
          "selected": false,
          "text": "All",
          "value": "$__all"
        },
        "datasource": "Prometheus-varnish",
        "definition": "label_values(varnish_up{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\", namespace=\"{{.Namespace}}\" }, pod)",
        "description": null,
        "error": null,
        "hide": 0,
        "includeAll": true,
        "label": "Pod",
        "multi": false,
        "name": "varnish_pod",
        "options": [],
        "query": {
          "query": "label_values(varnish_up{ service=\"{{.ServiceName}}\", namespace=\"{{.Namespace}}\", namespace=\"{{.Namespace}}\" }, pod)",
          "refId": "StandardVariableQuery"
        },
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 1,
        "tagValuesQuery": "",
        "tags": [],
        "tagsQuery": "",
        "type": "query",
        "useTags": false
      }
    ]
  },
  "time": {
    "from": "now-15m",
    "to": "now"
  },
  "timepicker": {
    "now": true,
    "nowDelay": "",
    "refresh_intervals": [
      "5s",
      "15s",
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
  "title": "{{.Title}}",
  "uid": "",
  "version": 1
}`
