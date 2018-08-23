#!/bin/sh

output_dir=$1

yamlblock=""
IFS=''
(kustomize build config/default; echo "---") | while read line; do
    if [[ "$line" = "---" ]]; then
        kind=$(echo "$yamlblock" | sed -n -e 's/^kind: //p')
        echo "$yamlblock" > $output_dir/varnishservice_$kind.yaml
        yamlblock=""
    elif [[ -z "$yamlblock" ]]; then
        yamlblock="$line"
    else
        yamlblock+="\n$line"
    fi
done
