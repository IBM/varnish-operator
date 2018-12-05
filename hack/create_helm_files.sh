#!/bin/bash

(kustomize build "config/default" | yq -c .) |
  while read -r json; do
    yq -y . <<< $json > "$1/manager_$(yq -r '.kind | ascii_downcase' <<< $json).yaml"
  done