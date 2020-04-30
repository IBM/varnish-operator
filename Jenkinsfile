#!groovy
@Library("icm-jenkins-common@0.117.0")
import com.ibm.icm.*

icmJenkinsProperties().
    rotateBuilds(numToKeep: 30, daysToKeep: 45).
    disableConcurrentBuilds().
    apply()

ibmCloud = [
    region: 'us-south',
    apiKeyId: 'icm_bluemix_1638245'] // id of Jenkins credential secret text

// for Docker
docker = [
    registry: 'us.icr.io',
    registryNamespace: 'icm-varnish',
    isLatest: false
]

releaseBranch = 'master'
docsBranch = 'gh-pages'
slackChannel = 'varnish-operator-notifications'
gitCredentialId = 'ApplicationId-icmautomation'
committerUsername = 'Core Engineering'
committerEmail = 'coreeng@us.ibm.com'

// for Helm chart publish
helmChart = 'varnish-operator'
artifactoryRoot = 'na.artifactory.swg-devops.com/artifactory'
artifactoryRepo = 'wcp-icm-helm-local'
artifactoryUserPasswordId = 'TAAS-Artifactory-User-Password-Global'

// For go modules download
goVirtualProxyRepo = 'wcp-icm-go-virtual'

dockerKeepDev = [
        releases: 0,
        days: 14,
        tagEndsWith: '-dev'
]

node('icm_agent_go') {
  GitInfo gitInfo = icmCheckoutStages(withTags: true) // By default the clone occurs without refs fetch

  runTests()

  def appVersion = icmGetAppVersion() // Get the version from version.txt
  boolean isTaggedCommit = icmGetTagsOnCommit().size() > 0 //used to avoid build runs on tag push
  boolean isNewRelease = !isTaggedCommit && gitInfo.branch == releaseBranch
  boolean isUntaggedCommit = !isTaggedCommit && gitInfo.branch != releaseBranch
  slack = icmSlackNotifier(slackChannel)

  HelmUtils helmUtils = new HelmUtils(this)

  // Builds for the release branch are triggered for any commit to it except tags (last commit shouldn't contain associated tags)
  // and if tag from version.txt exist in cloned repo jenkins completes the job silently with no errors
  // Builds are triggered for commits made on non-release branches except tags

  if (isNewRelease) {
    GitUtils gitUtils = new GitUtils(this)
    // If it's master and version.txt have changed - make a release (push a new git tag, docker images and helm chart)
    if (!gitUtils.doesTagExist(appVersion)) {
      println 'This is a release branch. Preparing the build...'
      try {
        icmCloudCliSetupStages(ibmCloud.apiKeyId, ibmCloud.region, CloudCliPluginConsts.CONTAINER_PLUGINS)
        dockerBuildPush(appVersion)
        helmUtils.setChartVersion(helmChart, appVersion)
        icmWithArtifactoryConfig(artifactoryRoot, artifactoryRepo, artifactoryUserPasswordId) {
          icmHelmChartPackagePublishStage(helmChart, it.config.createHelmPublish())
        }

        buildPushDocs()

        // Push the newly created release tag to origin
        gitUtils.setTag(appVersion, gitCredentialId)
        slack.info("Varnish Operator: $appVersion has been released. See the <https://github.ibm.com/TheWeatherCompany/icm-varnish-k8s-operator/releases/tag/${appVersion}|release notes>. For more debug info see the <${env.BUILD_URL}|Jenkins logs>")
      } catch (err) {
        String errorMessage = "Release $appVersion Failed! For more debug info see the <${env.BUILD_URL}|Jenkins logs>"
        slack.error(errorMessage)
        icmMarkBuildFailed(errorMessage)
      }
    }
    // If it's master and version.txt haven't changed - do nothing.

  // If it's a feature branch - push only artifacts (docker images, helm chart), but do not push a git tag
  } else if (isUntaggedCommit) {
    println 'This is a feature branch. Preparing the build...'
    icmCloudCliSetupStages(ibmCloud.apiKeyId, ibmCloud.region, CloudCliPluginConsts.CONTAINER_PLUGINS)
    appVersion = StringUtils.slugify("$appVersion-${gitInfo.branch}-dev", 128)
    dockerBuildPush(appVersion)
    helmUtils.setChartVersion(helmChart, appVersion)
    icmWithArtifactoryConfig(artifactoryRoot, artifactoryRepo, artifactoryUserPasswordId) {
      icmHelmChartPackagePublishStage(helmChart, it.config.createHelmPublish())
    }
  }
}

def dockerBuildPush(String appVersion) {
  stage('Docker build & push') {

    def stepsForParallel = [:]

    icmWithArtifactoryConfig(artifactoryRoot, goVirtualProxyRepo, artifactoryUserPasswordId) {
      stepsForParallel['Building varnish'] = {
        DockerImageInfo dockerImageInfo = new DockerImageInfo(docker.registry, docker.registryNamespace, 'varnish', appVersion, docker.isLatest)
        cleanupDockerImages(dockerImageInfo)
        icmDockerStages(dockerImageInfo, ['-f': 'Dockerfile.varnishd'])
      }
      stepsForParallel['Building varnish-metrics-exporter'] = {
        DockerImageInfo dockerImageInfo = new DockerImageInfo(docker.registry, docker.registryNamespace, 'varnish-metrics-exporter', appVersion, docker.isLatest)
        cleanupDockerImages(dockerImageInfo)
        icmDockerStages(dockerImageInfo, ['-f': 'Dockerfile.exporter'])
      }
      stepsForParallel['Building varnish-controller'] = {
        DockerImageInfo dockerImageInfo = new DockerImageInfo(docker.registry, docker.registryNamespace, 'varnish-controller', appVersion, docker.isLatest)
        cleanupDockerImages(dockerImageInfo)
        icmDockerStages(dockerImageInfo, ['-f': 'Dockerfile.controller', '--build-arg': "GOPROXY=\"https://${ARTIFACTORY_USER}:${ARTIFACTORY_PASS}@na.artifactory.swg-devops.com/artifactory/${ARTIFACTORY_REPO}/\""])
      }
      stepsForParallel['Building varnish-operator'] = {
        DockerImageInfo dockerImageInfo = new DockerImageInfo(docker.registry, docker.registryNamespace, 'varnish-operator', appVersion, docker.isLatest)
        cleanupDockerImages(dockerImageInfo)
        icmDockerStages(dockerImageInfo, ['-f': 'Dockerfile', '--build-arg': "GOPROXY=\"https://${ARTIFACTORY_USER}:${ARTIFACTORY_PASS}@na.artifactory.swg-devops.com/artifactory/${ARTIFACTORY_REPO}/\""])
      }

    parallel stepsForParallel
    }
  }
}

def buildPushDocs() {
  stage('Docs') {
    String url = steps.sh(returnStdout: true, script: "git remote get-url origin").replaceAll("https://", "").trim()
    sh(script: """
      gitbook install ./docs
      gitbook build ./docs docs_generated --log=debug --debug
      cd docs_generated/
      git init
      git add .
      git config --local user.name \"$committerUsername\"
      git config --local user.email \"$committerEmail\"
      git commit -m "Deploy Docs from Jenkins"
    """)

    UserPasswordByIdInfo userPasswordByIdInfo = new UserPasswordByIdInfo(this, gitCredentialId)
    userPasswordByIdInfo.withUserAndPassword { username, password ->
      sh(
        script: "git push --force --quiet \"https://$username:$password@$url\" HEAD:$docsBranch"
      )
    }

    sh("cd .. && rm -rf ./docs_generated")
  }
}

def runTests() {
  stage('Tests') {
    icmWithArtifactoryConfig(artifactoryRoot, goVirtualProxyRepo, artifactoryUserPasswordId) {
      sh(script: """
      export GO111MODULE=on
      export GOPROXY="https://${ARTIFACTORY_USER}:${ARTIFACTORY_PASS}@na.artifactory.swg-devops.com/artifactory/${ARTIFACTORY_REPO}/"
      go mod download
      golangci-lint run --timeout 2m --verbose
      go test ./pkg/... ./api/... -coverprofile=cover.out
      go tool cover -func=cover.out | tail -1 | awk '{print \"Total coverage: \" \$3}'
    """)
    }
  }
}

def cleanupDockerImages(DockerImageInfo dockerImageInfo) {
  icmDockerCleanupStage(dockerImageInfo, dockerKeepDev['tagEndsWith'], dockerKeepDev['releases'], dockerKeepDev['days'], true, false)
}
