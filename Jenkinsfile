#!groovy
@Library("icm-jenkins-common@0.36.1")
import com.ibm.icm.*

// for GitInfo
bxApiKeyId = "icm_bluemix_1638245"
region = "us-south"

// for Docker
dockerRegistry = "us.icr.io"
dockerRegistryNamespace = "icm-varnish"
varnishDockerImageName = "varnish"
operatorDockerImageName = "varnish-controller"
releaseBranch = "master"

// for Helm chart publish
helmChartPath = "varnish-operator"
artifactoryHostName = "na.artifactory.swg-devops.com/artifactory"
artifactoryRepo = "wcp-icm-helm-local"
artifactoryCredentialId = "TAAS-Artifactory-User-Password-Global"

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
    
    List<String> tags = icmGetTagsOnCommit()
    String repoVersion = new VersionUtils(this).getAppVersion()
    if (tags && tags.contains(repoVersion)) {
        stage("Helm Chart Publish") {
            sh "./hack/create_helm_files.sh ${helmChartPath}/templates"
            icmWithArtifactoryConfig(artifactoryHostName, artifactoryRepo, artifactoryCredentialId) { config, envNames, namesToValues ->
                icmHelmChartPackagePublish(helmChartPath, config.createHelmPublish())
            }
        }
    }
}
