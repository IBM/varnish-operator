{{/* vim: set filetype=mustache: */}}
{{/* Expand the name of the chart. */}}
{{- define "varnish-operator.name" -}}
  {{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "varnish-operator.fullname" -}}
  {{- if .Values.fullnameOverride -}}
    {{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
  {{- else -}}
    {{- $name := default .Chart.Name .Values.nameOverride -}}
    {{- if contains $name .Release.Name -}}
      {{- .Release.Name | trunc 63 | trimSuffix "-" -}}
    {{- else -}}
      {{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
    {{- end -}}
  {{- end -}}
{{- end -}}

{{/* Create chart name and version as used by the chart label. */}}
{{- define "varnish-operator.chart" -}}
  {{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/* use a release name that has a unique name per namespace */}}
{{- define "varnish-operator.releaseName" -}}
  {{- printf "%s-%s" .Release.Name .Values.namespace -}}
{{- end -}}

{{/* converts an array into a comma separated list */}}
{{- define "commaSeparatedList" -}}
  {{- range $index, $cmd := . -}}
    {{- if $index -}},{{- end -}}{{- $cmd -}}
  {{- end -}}
{{- end -}}
