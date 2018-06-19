#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

SCRIPT_ROOT=$(dirname ${BASH_SOURCE})/..
CODEGEN_PKG=${CODEGEN_PKG:-$(ls -d -1 $SCRIPT_ROOT/vendor/k8s.io/code-generator 2>/dev/null)}

$CODEGEN_PKG/generate-groups.sh deepcopy,defaulter,lister,informer \
icm-varnish-k8s-operator/operator/controller/pkg/client icm-varnish-k8s-operator/operator/controller/pkg/apis \
icm:v1alpha1