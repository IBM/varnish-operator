#!groovy
@Library("icm-jenkins-common@0.30.0")
import com.ibm.icm.*

region = 'us-south'
bxApiKeyId = 'icm_bluemix_1638245'
releaseBranch = 'master'
dockerRegistry = 'registry.ng.bluemix.net'
dockerRegistryNamespace = 'icm-varnish'
dockerImageName = 'varnish-controller'
varnishDockerImageName = 'varnish'
artifactoryHostName = "na.artifactory.swg-devops.com"
artifactoryRepo = "wcp-icm-helm-local"
artifactoryCredentialId='TAAS-Artifactory-User-Password-Global'

node('icm_slave') {
    sh "ln -s /etc/bluemix ~/.bluemix"
    GitInfo gitInfo = icmCheckoutStages()
    icmLoginToBx(bxApiKeyId, region, BxPluginConsts.CONTAINER_PLUGINS)
    DockerImageInfo operatorDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, dockerImageName,
            releaseBranch, gitInfo)
    icmDockerStages(operatorDockerImageInfo)
    DockerImageInfo varnishDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, varnishDockerImageName,
            releaseBranch, gitInfo)

    def buildOptions = ["-f":"Dockerfile.Varnish"]
    icmDockerStages(varnishDockerImageInfo, buildOptions)
    if (gitInfo.branch == releaseBranch) {
        sh './hack/create_helm_files.sh ./varnish-operator/templates'
        icmArtifactoryHelmChartPackageAndPublish('varnish-operator', artifactoryCredentialId, artifactoryHostName, artifactoryRepo)
    }
}