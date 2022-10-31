{{/*
Return the proper Varnish Operator image name
*/}}
{{- define "varnish-operator.image" -}}
    {{- if .Values.container.image }}
        {{- $image := .Values.container.image -}}
        {{- printf "%s" $image -}}
    {{- else -}}
        {{- $registryName := .Values.container.registry -}}
        {{- $repositoryName := .Values.container.repository -}}
        {{- $separator := ":" -}}
        {{- $termination := .Values.container.tag | toString -}}
        {{- if .Values.container.digest }}
            {{- $separator = "@" -}}
            {{- $termination = .Values.container.digest | toString -}}
        {{- end -}}
        {{- printf "%s/%s%s%s" $registryName $repositoryName $separator $termination -}}
    {{- end -}}
{{- end -}}
