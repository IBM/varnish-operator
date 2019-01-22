#!groovy
@Library("icm-jenkins-common@0.36.1")
import com.ibm.icm.*

region = "us-south"
bxApiKeyId = "icm_bluemix_1638245"
releaseBranch = "master"
dockerRegistry = "registry.ng.bluemix.net"
dockerRegistryNamespace = "icm-varnish"
operatorDockerImageName = "varnish-controller"
varnishDockerImageName = "varnish"
artifactoryHostName = "na.artifactory.swg-devops.com"
artifactoryRepo = "wcp-icm-helm-local"
artifactoryCredentialId="TAAS-Artifactory-User-Password-Global"

node("icm_slave") {
    sh "ln -s /etc/bluemix ~/.bluemix"
    GitInfo gitInfo = icmCheckoutStages()
    icmLoginToBx(bxApiKeyId, region, BxPluginConsts.CONTAINER_PLUGINS)

    DockerImageInfo varnishDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, varnishDockerImageName,
            releaseBranch, gitInfo)
    icmDockerStages(varnishDockerImageInfo, ["-f":"Dockerfile.Varnish"])

    DockerImageInfo operatorDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, operatorDockerImageName,
            releaseBranch, gitInfo)
    icmDockerStages(operatorDockerImageInfo)

    if (gitInfo.branch == releaseBranch) {
        sh "./hack/create_helm_files.sh ./varnish-operator/templates"
        icmWithArtifactoryConfig(artifactoryHostName, artifactoryRepo, artifactoryCredentialId) { config, envNames, namesToValues ->
            icmHelmChartPackagePublish("varnish-operator", config.createHelmPublish())
        }
    }
}
