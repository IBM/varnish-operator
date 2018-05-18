#!groovy
// @Library("icm-jenkins-common")
// import com.ibm.icm.*

// releaseBranch = 'master'
// dockerRegistryNamespace = 'icm-varnish'
// dockerImageName = 'varnish'

// defaultHelmChartInfo = new HelmChartInfo([
//         chart              : 'varnish',
//         releaseName        : 'varnish-cache',
// ])

// availableClusters =
//         [new DeployClusterInfo([
//                 region          : 'us-south',
//                 name            : 'icm-poc-shared',
//                 releaseNamespace: 'varnish-ns',
//                 bxApiKeyId      : 'icm_bluemix_1638245',
//                 dockerRegistry  : "registry.ng.bluemix.net"
//         ])]

// node {
//     GitInfo gitInfo = icmCheckoutStages()
//     icmInstallBxCliWithPlugins(BxPluginConsts.CONTAINER_PLUGINS)
//     availableClusters.each {
//         DockerImageInfo dockerImageInfo = icmGetDockerImageInfo(it.dockerRegistry, dockerRegistryNamespace,
//                 dockerImageName, releaseBranch, gitInfo)
//         icmDockerBuild(dockerImageInfo)
//         icmLoginToBx(it.bxApiKeyId, it.region, BxPluginConsts.CONTAINER_PLUGINS)

//         icmDockerPush(dockerImageInfo)
//         icmDeployWithHelmStages(it, defaultHelmChartInfo, dockerImageInfo)
//     }
// }