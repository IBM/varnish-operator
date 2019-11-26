#!groovy
@Library("icm-jenkins-common@0.94.0")
import com.ibm.icm.*

// for GitInfo
cloudApiKeyId = "icm_bluemix_1638245"
region = "us-south"

// for Docker
dockerRegistry = "us.icr.io"
dockerRegistryNamespace = "icm-varnish"
varnishDockerImageName = "varnish"
varnishControllerDockerImageName = "varnish-controller"
varnishMetricsExporterDockerImageName = "varnish-metrics-exporter"
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
    DockerImageInfo varnishControllerDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, varnishControllerDockerImageName,
            releaseBranch, gitInfo)
    DockerImageInfo varnishMetricsExporterDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, varnishMetricsExporterDockerImageName,
            releaseBranch, gitInfo)

    DockerImageInfo operatorDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, operatorDockerImageName,
            releaseBranch, gitInfo)

    String repoVersion = versionUtils.getAppVersion()

    if ( branch == releaseBranch ) {

        if ( ! gitUtils.doesTagExist(repoVersion) ) {

            icmDockerStages(varnishDockerImageInfo, ["-f":"Dockerfile.varnishd"])
            icmDockerStages(varnishMetricsExporterDockerImageInfo, ["-f":"Dockerfile.exporter"])
            icmDockerStages(varnishControllerDockerImageInfo, ["-f":"Dockerfile.controller"])
            icmDockerStages(operatorDockerImageInfo)

            icmWithArtifactoryConfig(artifactoryRoot, artifactoryRepo, artifactoryUserPasswordId) {
                icmHelmChartPackagePublishStage(helmChart, it.config.createHelmPublish())
            }

            gitUtils.setTag(repoVersion, credentialsId, 'origin')
        }
    } else {

        icmDockerStages(varnishDockerImageInfo, ["-f":"Dockerfile.varnishd"])
        icmDockerStages(varnishMetricsExporterDockerImageInfo, ["-f":"Dockerfile.exporter"])
        icmDockerStages(varnishControllerDockerImageInfo, ["-f":"Dockerfile.controller"])
        icmDockerStages(operatorDockerImageInfo)

        stage('Set the helm chart version') {
            helmUtils.setChartVersion(helmChart, repoVersion + '-' + gitInfo.branch)
        }
        icmWithArtifactoryConfig(artifactoryRoot, artifactoryRepo, artifactoryUserPasswordId) {
            icmHelmChartPackagePublishStage(helmChart, it.config.createHelmPublish())
        }
    }
}
