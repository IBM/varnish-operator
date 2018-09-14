#!groovy
@Library("icm-jenkins-common@0.19.0")
import com.ibm.icm.*

region = 'us-south'
bxApiKeyId = 'icm_bluemix_1638245'
releaseBranch = 'master'
dockerRegistry = 'registry.ng.bluemix.net'
dockerRegistryNamespace = 'icm-varnish'
dockerImageName = 'varnish-controller'
artifactoryHostName = "na.artifactory.swg-devops.com"
artifactoryRepo = "wcp-icm-helm-local"
artifactoryCredentialId='TAAS-Artifactory-User-Password-Global'

node {
    GitInfo gitInfo = icmCheckoutStages()
    icmDockerBuildStage(gitInfo)
    icmInstallBxCliWithPluginsStage(BxPluginConsts.CONTAINER_PLUGINS)
    DockerImageInfo dockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, dockerImageName, releaseBranch, gitInfo)
    icmLoginToBxStage(bxApiKeyId, region, BxPluginConsts.CONTAINER_PLUGINS)
    icmDockerPushStage(dockerImageInfo, gitInfo)
    if (gitInfo.branch == releaseBranch) {
        icmArtifactoryHelmChartPackageAndPublish('varnish-operator', artifactoryCredentialId, artifactoryHostName, artifactoryRepo)
    }
}