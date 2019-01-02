# Quick Start for Varnish Operator

## Prerequisites

* Kubernetes v1.11 or newer

## If Helm Is Not Configured with Artifactory Containing Helm Charts

1. Get access to [our Artifactory](https://na.artifactory.swg-devops.com) (TODO: add more instructions)
1. Generate an API Key on the Artifactory website
1. Follow [these instructions](https://www.jfrog.com/confluence/display/RTF/Helm+Chart+Repositories), where username is your email, password is your API Key:
    ```sh
    helm repo add wcp-icm-helm-virtual https://na.artifactory.swg-devops.com/artifactory/wcp-icm-helm-virtual --username=<your-email> --password=<api-key>
    helm repo update
    ```

## If Cluster Is Not Configured to Pull Docker Containers from IBM Container Registry

1. Get access to our IBM Cloud account (TODO: add more instructions)
1. use CLI to generate token:
    ```sh
    ibmcloud cr token-add --non-expiring --description 'for Varnish operator'
    ```
   And save the token
1. Create namespace in which varnish operator will be deployed
    ```sh
    kubectl create namespace <namespace>
    ```
1. Add token to that namespace:
    ```sh
    kubectl create secret docker-registry docker-reg-secret --namespace <namespace> --docker-server=registry.ng.bluemix.net --docker-username=token --docker-password=<token> --docker-email=dummy@ignore.me
    ```
   (`docker-reg-secret` used here as name of secret, but any name can be used)

## Install Varnish Operator

1. Create override `values.yaml` to specify namespace which has container registry token:
    ```yaml
    namespace: <namespace> # fill me in
    controllerImage:
      imagePullSecretName: docker-reg-secret # or whatever name you gave to the secret
    ```
1. Install varnish operator:
    ```sh
    helm upgrade --install varnish-operator wcp-icm-helm-virtual/varnish-operator --wait --namespace <namespace>
    ```

## Create a VarnishService

1. Add same IBM Container Registry token as above to target namespace where deployment that needs Varnish exists:
    ```sh
    kubectl create secret docker-registry docker-reg-secret --namespace <target-namespace> --docker-server=registry.ng.bluemix.net --docker-username=token --docker-password=<token> --docker-email=dummy@ignore.me
    ```
    (`docker-reg-secret` used here as name of secret, but any name can be used)
1. Create VarnishService yaml for spec of Varnishes:
    ```yaml
    apiVersion: icm.ibm.com/v1alpha1
    kind: VarnishService
    metadata:
      name: varnishservice-example
      namespace: <target-namespace> # fill me in
    spec:
      deployment:
        replicas: 3
        imagePullSecretName: docker-reg-secret # or whatever name you gave to the secret
      service:
        selector:
          app: MyApp # replace with selector for deployment that needs Varnish
        ports:
          - port: 8080 # or whatever port to expose through service
            targetPort: 8080 # or whatever port exposed on target deployment (optional if same as `port`)
    ```
1. Apply VarnishService yaml:
    ```sh
    kubectl apply -f <varnish-service>.yaml
    ```

## What You Should See

After running `kubectl apply`, you should see:

* Deployment called `<varnish-service-name>-deployment`
* Service called `<varnish-service-name>-cached` which uses Varnish for caching
* Service called `<varnish-service-name>-nocached` which bypasses Varnish
* ConfigMap called `vcl-file` (or whatever name specified in VarnishService config) containing VCL files that Varnish is using
* Role/Rolebinding/ClusterRole/ClusterRoleBinding/ServiceAccount combination for permissions

## What To Do Next

It is now possible to edit the ConfigMap to customize the VCL files. All files in the ConfigMap will be loaded into Varnish, but note that at least `default.vcl` and `backends.vcl.tmpl` must exist.

See the [README](README.md) for more detailed information about setup, like additional parameters for varnish-operator Helm chart overrides or VarnishService yaml.