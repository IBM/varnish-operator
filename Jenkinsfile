#!groovy
@Library("icm-jenkins-common@0.91.0")
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
credentialsId = 'ApplicationId-icmautomation'

// for Helm chart publish
helmChart = "varnish-operator"
artifactoryRoot = "na.artifactory.swg-devops.com/artifactory"
artifactoryRepo = "wcp-icm-helm-local"
artifactoryUserPasswordId = "TAAS-Artifactory-User-Password-Global"

VersionUtils versionUtils = new VersionUtils(this)
GitUtils gitUtils = new GitUtils(this)
HelmUtils helmUtils = new HelmUtils(this)

node("icm_slave") {
    GitInfo gitInfo = icmCheckoutStages()
    String branch = gitInfo.branch
    gitInfo.branch = StringUtils.slugify(gitInfo.branch)

    icmLoginToCloud(cloudApiKeyId, region, CloudCliPluginConsts.CONTAINER_PLUGINS)

    DockerImageInfo varnishDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, varnishDockerImageName,
            releaseBranch, gitInfo)

    DockerImageInfo operatorDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, operatorDockerImageName,
            releaseBranch, gitInfo)

    String repoVersion = versionUtils.getAppVersion()

    if ( branch == releaseBranch ) {

        if ( ! gitUtils.doesTagExist(repoVersion) ) {

            icmDockerStages(varnishDockerImageInfo, ["-f":"Dockerfile.Varnish"])
            icmDockerStages(operatorDockerImageInfo)

            icmWithArtifactoryConfig(artifactoryRoot, artifactoryRepo, artifactoryUserPasswordId) {
                icmHelmChartPackagePublishStage(helmChart, it.config.createHelmPublish())
            }

            gitUtils.setTag(repoVersion, credentialsId, 'origin')
        }
    } else {
        icmDockerStages(varnishDockerImageInfo, ["-f":"Dockerfile.Varnish"])
        icmDockerStages(operatorDockerImageInfo)
        stage('Set the helm chart version') {
            helmUtils.setChartVersion(helmChart, repoVersion + '-' + gitInfo.branch)
        }
        icmWithArtifactoryConfig(artifactoryRoot, artifactoryRepo, artifactoryUserPasswordId) {
            icmHelmChartPackagePublishStage(helmChart, it.config.createHelmPublish())
        }
    }
}
