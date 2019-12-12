#!groovy
@Library("icm-jenkins-common@0.96.0")
import com.ibm.icm.*

icmJenkinsProperties().
    rotateBuilds(numToKeep: 30, daysToKeep: 45).
    disableConcurrentBuilds().
    apply()

ibmCloud = [
    region: 'us-south',
    apiKeyId: 'icm_bluemix_1638245'] // id of Jenkins credential secret text

// for Docker
dockerRegistry = "us.icr.io"
dockerRegistryNamespace = "icm-varnish"
varnishDockerImageName = "varnish"
varnishControllerDockerImageName = "varnish-controller"
varnishMetricsExporterDockerImageName = "varnish-metrics-exporter"
operatorDockerImageName = "varnish-operator"
releaseBranch = "master"
slackChannel = 'varnish-operator-notifications'
gitCredentialId = 'ApplicationId-icmautomation'

// for Helm chart publish
helmChart = "varnish-operator"
artifactoryRoot = "na.artifactory.swg-devops.com/artifactory"
artifactoryRepo = "wcp-icm-helm-local"
artifactoryUserPasswordId = "TAAS-Artifactory-User-Password-Global"

VersionUtils versionUtils = new VersionUtils(this)
HelmUtils helmUtils = new HelmUtils(this)

node("icm_slave") {
    GitInfo gitInfo = icmCheckoutStages(withTags: true) // By default the clone occurs without refs fetch
    String branch = gitInfo.branch
    gitInfo.branch = StringUtils.slugify(gitInfo.branch)

    DockerImageInfo varnishDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, varnishDockerImageName,
            releaseBranch, gitInfo)
    DockerImageInfo varnishControllerDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, varnishControllerDockerImageName,
            releaseBranch, gitInfo)
    DockerImageInfo varnishMetricsExporterDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, varnishMetricsExporterDockerImageName,
            releaseBranch, gitInfo)
    DockerImageInfo operatorDockerImageInfo = icmGetDockerImageInfo(dockerRegistry, dockerRegistryNamespace, operatorDockerImageName,
            releaseBranch, gitInfo)

    def appVersion = icmGetAppVersion() // Get the version from version.txt
    boolean isTaggedCommit = icmGetTagsOnCommit().size() > 0 //used to avoid build runs on tag push
    boolean isNewRelease = !isTaggedCommit && gitInfo.branch == releaseBranch
    boolean isUntaggedCommit = !isTaggedCommit && gitInfo.branch != releaseBranch
    slack = icmSlackNotifier(slackChannel)

    if (isNewRelease) {
        GitUtils gitUtils = new GitUtils(this)
        if (!gitUtils.doesTagExist(appVersion)) {
            println 'This is a release branch. Preparing the build...'
            try {
                icmCloudCliSetupStages(ibmCloud.apiKeyId, ibmCloud.region, CloudCliPluginConsts.CONTAINER_PLUGINS)
                icmDockerStages(varnishDockerImageInfo, ["-f": "Dockerfile.varnishd"])
                icmDockerStages(varnishMetricsExporterDockerImageInfo, ["-f": "Dockerfile.exporter"])
                icmDockerStages(varnishControllerDockerImageInfo, ["-f": "Dockerfile.controller"])
                icmDockerStages(operatorDockerImageInfo)
                helmUtils.setChartVersion(helmChart, appVersion)
                icmWithArtifactoryConfig(artifactoryRoot, artifactoryRepo, artifactoryUserPasswordId) {
                    icmHelmChartPackagePublishStage(helmChart, it.config.createHelmPublish())
                }
                // Push the newly created release tag to origin
                gitUtils.setTag(appVersion, gitCredentialId)
                slack.info("icm-varnish-operator: $appVersion has been released. See the <https://github.ibm.com/TheWeatherCompany/icm-cassandra/releases/tag/${appVersion}|release notes>. For more debug info see the <${env.BUILD_URL}|Jenkins logs>")
            } catch (err) {
                String errorMessage = "Release $appVersion Failed! For more debug info see the <${env.BUILD_URL}|Jenkins logs>"
                slack.error(errorMessage)
                icmMarkBuildFailed(errorMessage)
            }
        }
    } else if (isUntaggedCommit) {
        println 'This is a feature branch. Preparing the build...'
        icmCloudCliSetupStages(ibmCloud.apiKeyId, ibmCloud.region, CloudCliPluginConsts.CONTAINER_PLUGINS)
        appVersion += '-' + StringUtils.slugify(gitInfo.branch)
        icmDockerStages(varnishDockerImageInfo, ["-f":"Dockerfile.varnishd"])
        icmDockerStages(varnishMetricsExporterDockerImageInfo, ["-f":"Dockerfile.exporter"])
        icmDockerStages(varnishControllerDockerImageInfo, ["-f":"Dockerfile.controller"])
        icmDockerStages(operatorDockerImageInfo)
        helmUtils.setChartVersion(helmChart, appVersion)
        icmWithArtifactoryConfig(artifactoryRoot, artifactoryRepo, artifactoryUserPasswordId) {
            icmHelmChartPackagePublishStage(helmChart, it.config.createHelmPublish())
        }
    }
}
