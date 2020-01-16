#!groovy
@Library("icm-jenkins-common@0.100.0")
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

node('icm_slave_go') {
  GitInfo gitInfo = icmCheckoutStages(withTags: true) // By default the clone occurs without refs fetch

  stage('Tests') {
    sh(script: """
      export GO111MODULE=on
      go mod download
      golangci-lint run
      go test ./pkg/... ./api/... -coverprofile=cover.out
      go tool cover -func=cover.out | tail -1 | awk '{print \"Total coverage: \" \$3}'
    """)
  }

  def appVersion = icmGetAppVersion() // Get the version from version.txt
  boolean isTaggedCommit = icmGetTagsOnCommit().size() > 0 //used to avoid build runs on tag push
  boolean isNewRelease = !isTaggedCommit && gitInfo.branch == releaseBranch
  boolean isUntaggedCommit = !isTaggedCommit && gitInfo.branch != releaseBranch
  slack = icmSlackNotifier(slackChannel)

  HelmUtils helmUtils = new HelmUtils(this)

  if (isNewRelease) {
    GitUtils gitUtils = new GitUtils(this)
    if (!gitUtils.doesTagExist(appVersion)) {
      println 'This is a release branch. Preparing the build...'
      try {
        icmCloudCliSetupStages(ibmCloud.apiKeyId, ibmCloud.region, CloudCliPluginConsts.CONTAINER_PLUGINS)
        def buildRes = dockerBuildPush(appVersion, gitInfo) //use var assignment here to avoid `unable to resolve class Gitinfo`
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
  } else if (isUntaggedCommit) {
    println 'This is a feature branch. Preparing the build...'
    icmCloudCliSetupStages(ibmCloud.apiKeyId, ibmCloud.region, CloudCliPluginConsts.CONTAINER_PLUGINS)
    appVersion += '-' + StringUtils.slugify(gitInfo.branch)
    def buildRes = dockerBuildPush(appVersion, gitInfo) //use var assignment here to avoid `unable to resolve class Gitinfo`
    helmUtils.setChartVersion(helmChart, appVersion)
    icmWithArtifactoryConfig(artifactoryRoot, artifactoryRepo, artifactoryUserPasswordId) {
      icmHelmChartPackagePublishStage(helmChart, it.config.createHelmPublish())
    }
  }
}

def dockerBuildPush(String appVersion, GitInfo gitinfo) {
  stage('Docker build & push') {

    def stepsForParallel = [:]

    stepsForParallel['Building varnish'] = {
      icmDockerStages(new DockerImageInfo(docker.registry, docker.registryNamespace, 'varnish', appVersion, docker.isLatest), ['-f': 'Dockerfile.varnishd'])
    }
    stepsForParallel['Building varnish-metrics-exporter'] = {
      icmDockerStages(new DockerImageInfo(docker.registry, docker.registryNamespace, 'varnish-metrics-exporter', appVersion, docker.isLatest), ['-f': 'Dockerfile.exporter'])
    }
    stepsForParallel['Building varnish-controller'] = {
      icmDockerStages(new DockerImageInfo(docker.registry, docker.registryNamespace, 'varnish-controller', appVersion, docker.isLatest), ['-f': 'Dockerfile.controller'])
    }
    stepsForParallel['Building varnish-operator'] = {
      icmDockerStages(new DockerImageInfo(docker.registry, docker.registryNamespace, 'varnish-operator', appVersion, docker.isLatest))
    }

    parallel stepsForParallel
  }
}

def buildPushDocs() {
  stage('Docs') {
    String url = steps.sh(returnStdout: true, script: "git remote get-url origin").replaceAll("https://", "").trim()
    sh(script: """
      gitbook install
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
