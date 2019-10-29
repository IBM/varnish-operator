#!groovy
@Library("icm-jenkins-common@0.84.0")
import com.ibm.icm.*

// for GitInfo
cloudApiKeyId = "icm_bluemix_1638245"
region = "us-south"

// for Docker
dockerRegistry = "us.icr.io"
dockerRegistryNamespace = "icm-varnish"
varnishDockerImageName = "varnish"
operatorDockerImageName = "varnish-operator"
releaseBranch = "master"

// for Helm chart publish
helmChart = "varnish-operator"
artifactoryRoot = "na.artifactory.swg-devops.com/artifactory"
artifactoryRepo = "wcp-icm-helm-local"
artifactoryUserPasswordId = "TAAS-Artifactory-User-Password-Global"

node("icm_slave") {
    GitInfo gitInfo = icmCheckoutStages()
    icmLoginToCloud(cloudApiKeyId, region, CloudCliPluginConsts.CONTAINER_PLUGINS)

    DockerImageInfo varnishDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, varnishDockerImageName,
            releaseBranch, gitInfo)
    icmDockerStages(varnishDockerImageInfo, ["-f": "Dockerfile.Varnish"])

    DockerImageInfo operatorDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, operatorDockerImageName,
            releaseBranch, gitInfo)
    icmDockerStages(operatorDockerImageInfo)

    List<String> tags = icmGetTagsOnCommit()
    String repoVersion = new VersionUtils(this).getAppVersion()
    if (tags && tags.contains(repoVersion)) {
        stage("Helm Chart Publish") {
            icmWithArtifactoryConfig(artifactoryRoot, artifactoryRepo, artifactoryUserPasswordId) {
                icmHelmChartPackagePublish(helmChart, it.config.createHelmPublish())
            }
        }
    }
}
