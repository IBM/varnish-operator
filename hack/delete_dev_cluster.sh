#!/bin/bash

cluster_name="e2e-tests"

function usage {
  cat << !
USAGE: $0 [-c cluster]

Deletes a dev cluster

-c|--cluster   | cluster
!
}

while (( "$#" )); do
  opt="$1"; shift;
  case "$opt" in
    "-c"|"--cluster") cluster_name="$1" shift;;
    *) echo "invalid option: \""$opt"\"" >&2; usage; exit 1;;
  esac
done

kind delete cluster --name $cluster_name > /dev/null 2>&1

if [[ -f ./e2e-tests-kubeconfig ]]; then
 rm -f ./e2e-tests-kubeconfig
fi
