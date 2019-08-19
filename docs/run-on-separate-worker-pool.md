# Running Varnish pods on a separate IKS worker pool

This example shows how to create an IKS worker pool and make Varnish pods run strictly on its workers, one per node.

References:
 * [How to create IKS clusters and worker pools.](https://console.bluemix.net/docs/containers/cs_clusters.html#clusters)
 * [Taints and Tolerations](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/)
 * [Affinity and anti-affinity](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity)
 
1. Create a worker pool in your cluster assuming you already have a cluster called `test-cluster`

    ```bash
    $ #Find out the available zones to your cluster
    $ ibmcloud ks cluster-get --cluster test-cluster | grep "Worker Zones" # Get the 
    Worker Zones:           dal10
    $ #Find out what machine type are available in your zone  
    $ ibmcloud ks machine-types --zone dal10
    OK
    Name                      Cores   Memory   Network Speed   OS             Server Type   Storage   Secondary Storage   Trustable   
    u2c.2x4                   2       4GB      1000Mbps        UBUNTU_16_64   virtual       25GB      100GB               false   
    ms2c.4x32.1.9tb.ssd       4       32GB     10000Mbps       UBUNTU_16_64   physical      2000GB    960GB               false   
    ms2c.16x64.1.9tb.ssd      16      64GB     10000Mbps       UBUNTU_16_64   physical      2000GB    960GB               true   
    ms2c.28x256.3.8tb.ssd     28      256GB    10000Mbps       UBUNTU_16_64   physical      2000GB    1920GB              true   
       ...
    $ #Create a worker pool. 
    $ ibmcloud ks worker-pool-create --name varnish-worker-pool --cluster test-cluster --machine-type u2c.2x4 --size-per-zone 2 --hardware shared
    OK 
    $ #Verify your worker pool is created
    $ ibmcloud ks worker-pools --cluster test-cluster
    Name                  ID                                         Machine Type          Workers   
    default               91ed9433e7bf4dc7b8348ae1022f9f27-89d7d12   b2c.16x64.encrypted   3   
    varnish-worker-pool   91ed9433e7bf4dc7b8348ae1022f9f27-c5b13da   u2c.2x4.encrypted     2   
    OK
    $ #Add your zone to your worker pool. First, find out your VLAN IDs
    $ ibmcloud ks vlans --zone dal10
    OK
    ID        Name   Number   Type      Router         Supports Virtual Workers   
    2315193          1690     private   bcr02a.dal10   true   
    2315191          1425     public    fcr02a.dal10   true
    $ #Use the VLAN IDs above to add your zone to the worker pool
    $ ibmcloud ks zone-add --zone dal10 --cluster test-cluster --worker-pools varnish-worker-pool --private-vlan 2315193 --public-vlan 2315191
    OK
    $ #Verify that worker nodes provision in the zone that you've added
    $ ibmcloud ks workers --cluster test-cluster --worker-pool varnish-worker-pool
    OK
    ID                                                  Public IP   Private IP   Machine Type        State               Status                          Zone    Version   
    kube-dal10-cr91ed9433e7bf4dc7b8348ae1022f9f27-w58   -           -            u2c.2x4.encrypted   provision_pending   Preparing to provision worker   dal10   1.11.7_1543   
    kube-dal10-cr91ed9433e7bf4dc7b8348ae1022f9f27-w59   -           -            u2c.2x4.encrypted   provision_pending   -                               dal10   1.11.7_1543   
    ```
    
    Wait until your worker pool nodes change their state to `normal` and status to `Ready`.
    
    ```bash
    $ ibmcloud ks workers --cluster test-cluster --worker-pool varnish-worker-pool
    OK
    ID                                                  Public IP       Private IP      Machine Type        State    Status   Zone    Version   
    kube-dal10-cr91ed9433e7bf4dc7b8348ae1022f9f27-w58   169.61.218.68   10.94.177.179   u2c.2x4.encrypted   normal   Ready    dal10   1.11.7_1543   
    kube-dal10-cr91ed9433e7bf4dc7b8348ae1022f9f27-w59   169.61.218.94   10.94.177.180   u2c.2x4.encrypted   normal   Ready    dal10   1.11.7_1543
    ```
    
1. Taint created nodes to repel pods that don't have required toleration. 

    ```bash
    $ #Setup kubectl
    $ ibmcloud ks cluster-config --cluster test-cluster 
    OK
    The configuration for test-cluster was downloaded successfully.
    
    Export environment variables to start using Kubernetes.
    
    export KUBECONFIG=/home/me/.bluemix/plugins/container-service/clusters/test-cluster/kube-config-dal10-test-cluster.yml
    
    $ export KUBECONFIG=/home/me/.bluemix/plugins/container-service/clusters/test-cluster/kube-config-dal10-test-cluster.yml
    $ #Find your nodes using kubectl. First get your worker pool ID and then use it to select your nodes
    $ ibmcloud ks worker-pools --cluster test-cluster 
    Name                  ID                                         Machine Type          Workers   
    default               91ed9433e7bf4dc7b8348ae1022f9f27-89d7d12   b2c.16x64.encrypted   3   
    varnish-worker-pool   91ed9433e7bf4dc7b8348ae1022f9f27-c5b13da   u2c.2x4.encrypted     2   
    $ kubectl get nodes --selector ibm-cloud.kubernetes.io/worker-pool-id=91ed9433e7bf4dc7b8348ae1022f9f27-c5b13da
    NAME            STATUS   ROLES    AGE   VERSION
    10.94.177.179   Ready    <none>   16m   v1.11.7+IKS
    10.94.177.180   Ready    <none>   15m   v1.11.7+IKS
    $ #Taint those nodes
    $ kubectl taint node 10.94.177.179 role=varnish:NoSchedule #Do not schedule here not Varnish pods
    node/10.94.177.179 tainted
    $ kubectl taint node 10.94.177.179 role=varnish:NoExecute #Evict not Varnish pods if they already running here
    node/10.94.177.179 tainted
    $ kubectl taint node 10.94.177.180 role=varnish:NoSchedule #Do not schedule here not Varnish pods
    node/10.94.177.180 tainted
    $ kubectl taint node 10.94.177.180 role=varnish:NoExecute #Evict not Varnish pods if they already running here
    node/10.94.177.180 tainted
    ```
    
    This prevents all pods from scheduling on that node unless you already have pods with matching toleration
    
1. Label the nodes for the ability to schedule your varnish pods only on that nodes. Those labels will be used in your VarnishService configuration later.

    ```bash
    $ kubectl label node 10.94.177.179 role=varnish-cache
    node/10.94.177.179 labeled
    $ kubectl label node 10.94.177.180 role=varnish-cache
    node/10.94.177.180 labeled 
    ```
1. Define your VarnishService spec with necessary affinity and toleration configuration

    4.1 Define pods anti-affinity to not co-locate replicas on a node.
    
    ```yaml
    metadata:
      labels:
        role: varnish-cache
    spec:
      statefulSet:
        affinity:
          podAntiAffinity:
            requiredDuringSchedulingIgnoredDuringExecution:
              - labelSelector:
                  matchExpressions:
                    - key: role
                      operator: In
                      values:
                        - varnish-cache
                topologyKey: "kubernetes.io/hostname"
    ```
    That will make sure that two varnish pods doesn't get scheduled on one node. Kubernetes makes the decision based on labels we've set in the spec
    
    4.2 Define pods node affinity
    
    ```yaml
    spec:
      statefulSet:
        affinity:
          nodeAffinity:
            requiredDuringSchedulingIgnoredDuringExecution:
              nodeSelectorTerms:
              - matchExpressions:
                - key: role
                  operator: In
                  values:
                    - varnish-cache
    ```
    That will make kubernetes schedule varnish pods only on our worker pool nodes. The labels used here are the ones we've assigned to the node in step 3
    
    4.3 Define pods tolerations
    
    ```yaml
    spec:
      statefulSet:
        tolerations:
          - key: "role"
            operator: "Equal"
            value: "varnish"
            effect: "NoSchedule"
          - key: "role"
            operator: "Equal"
            value: "varnish"
            effect: "NoExecute"
    ```
    In step 2 we made our node repel all pods that don't have specific tolerations. Here we added those tolerations to be eligible for scheduling on those nodes. The values are the ones we used when tainted our nodes in step 2. 
    
5. Apply your configuration.

    This step assumes you have varnish operator [installed](#installation) and the namespace has the necessary secret [installed](#configuring-access).
    
    Complete VarnishService configuration example:
    
    ```yaml
    apiVersion: icm.ibm.com/v1alpha1
    kind: VarnishService
    metadata:
      labels:
        role: varnish-cache
      name: varnish-in-worker-pool
      namespace: varnish-ns
    spec:
      vclConfigMap:
        name: varnish-worker-pool-files
        backendsFile: backends.vcl
        defaultFile: default.vcl
      statefulSet:
        replicas: 2
        container:
          resources:
            limits:
              cpu: 100m
              memory: 256Mi
            requests:
              cpu: 100m
              memory: 256Mi
          readinessProbe:
            exec:
              command: [/usr/bin/varnishadm, ping]
          imagePullSecret: docker-reg-secret
        affinity:
          podAntiAffinity:
            requiredDuringSchedulingIgnoredDuringExecution:
              - labelSelector:
                  matchExpressions:
                    - key: role
                      operator: In
                      values:
                        - varnish-cache
                topologyKey: "kubernetes.io/hostname"
          nodeAffinity:
            requiredDuringSchedulingIgnoredDuringExecution:
              nodeSelectorTerms:
              - matchExpressions:
                - key: role
                  operator: In
                  values:
                    - varnish-cache
        tolerations:
          - key: "role"
            operator: "Equal"
            value: "varnish-cache"
            effect: "NoSchedule"
          - key: "role"
            operator: "Equal"
            value: "varnish-cache"
            effect: "NoExecute"
      service:
        selector:
          app: HttPerf
        varnishPort:
          name: varnish
          port: 2035
          targetPort: 8080
        varnishExporterPort:
          name: varnishexporter
          port: 9131
    ```
    Apply your configuration:
    ```bash
    $ kubectl apply -f varnish-in-worker-pool.yaml
    varnishservice.icm.ibm.com/varnish-in-worker-pool created
    ```
    Here the operator will create all pods with specified configuration
6. See your pods being scheduled strictly on your worker pool and spread across different nodes.
    ```bash
    $ kubectl get pods -n varnish-ns -o wide --selector role=varnish-cache
    NAME                                                         READY   STATUS    RESTARTS   AGE   IP               NODE            NOMINATED NODE
    varnish-in-worker-pool-varnish-statefulset-0                 1/1     Running   0          6m    172.30.244.65    10.94.177.179   <none>
    varnish-in-worker-pool-varnish-statefulset-1                 1/1     Running   0          6m    172.30.136.129   10.94.177.180   <none>

    ```
    Check the `NODE` column. The value will be different for each pod.
    
    Note that you won't be able to run more pods than you have nodes. The anti-affinity rule will not allow two pods being co-located on one node.
    This behaviour can be changed by using an anti-affinity type called `preferredDuringSchedulingIgnoredDuringExecution`: 
    
    ```yaml
    spec:
      statefulSet:
        affinity:
          podAntiAffinity:
            preferredDuringSchedulingIgnoredDuringExecution:
            - weight: 1
              podAffinityTerm:
                labelSelector:
                  matchExpressions:
                  - key: role
                    operator: In
                    values:
                    - varnish-cache
                topologyKey: "kubernetes.io/hostname"
    ```
     It will still ask Kubernetes to spread pods onto different nodes but also allow to co-locate them if there are more pods than nodes.
     
    Also keep in mind that in such configuration the pods can be scheduled to your worker pool only. If the worker pool is deleted the pods will hang in `Pending` state until new nodes with the same configuration are added to the cluster.
