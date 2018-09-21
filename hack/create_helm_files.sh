#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..

output_dir=$1

yamlblock=""
IFS=''
(kustomize build "$SCRIPT_ROOT/config/default"; echo -e "---") | while read line; do
    if [[ "$line" = "---" ]]; then
        kind=$(echo -e "$yamlblock" | sed -n -e 's/^kind: //p')
        if [[ "$kind" = "ClusterRoleBinding" ]]; then
            yamlblock=$(echo -e "$yamlblock" | sed 's/namespace: .*/namespace: \{\{ .Values.namespace | quote \}\}/')
        fi
        echo -e "$yamlblock" > $output_dir/varnishservice_$kind.yaml
        yamlblock=""
    elif [[ -z "$yamlblock" ]]; then
        yamlblock="$line"
    else
        yamlblock+="\n$line"
    fi
done
