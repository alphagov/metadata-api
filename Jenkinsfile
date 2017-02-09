#!/usr/bin/env groovy

REPOSITORY = 'metadata-api'

node {
  def govuk = load '/var/lib/jenkins/groovy_scripts/govuk_jenkinslib.groovy'

  try {
    stage("Checkout") {
       checkout scm
    }

    stage("Setup environment") {
      govuk.setEnvar("REPO", "alphagov/metadata-api")
      govuk.setEnvar("GOPATH", "${PWD}/gopath")
      govuk.setEnvar("GO_GITHUB_PATH", "${GOPATH}/src/github.com")
      govuk.setEnvar("BUILD_PATH", "${GO_GITHUB_PATH}/${REPO}")
    }

    wrap([$class: 'AnsiColorBuildWrapper']) {
      stage("Test") {
        // TODO split test from build in Makefile
      }

      stage("Build") {
        sshagent(['govuk-ci-ssh-key']) {
          sh "rm -rf ${GOPATH} && mkdir -p ${GOPATH}/bin ${BUILD_PATH}"
          sh "rsync -a ./ ${BUILD_PATH} --exclude=gopath"
          sh "cd ${BUILD_PATH} && make"
        }
      }
    }

    stage("Archive artefact") {
      archiveArtifacts "metadata-api"
    }

    if (env.BRANCH_NAME == 'master') {
      stage("Push release tag") {
        govuk.pushTag(REPOSITORY, env.BRANCH_NAME, 'release_' + env.BUILD_NUMBER)
      }

      stage("Deploy on Integration") {
        govuk.deployIntegration(REPOSITORY, env.BRANCH_NAME, 'release', 'deploy')
      }
    }

  } catch (e) {
      currentBuild.result = "FAILED"
      step([$class: 'Mailer',
            notifyEveryUnstableBuild: true,
            recipients: 'govuk-ci-notifications@digital.cabinet-office.gov.uk',
            sendToIndividuals: true])
      throw e
  }
}
