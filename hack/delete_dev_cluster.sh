#!/bin/bash
cluster_name="e2e-tests"

kind delete cluster --name $cluster_name > /dev/null 2>&1
rm -f ./e2e-tests-kubeconfig