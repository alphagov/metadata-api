#!/usr/bin/env groovy

REPOSITORY = 'metadata-api'

node {
  env.REPO      = 'alphagov/metadata-api'
  env.BUILD_DIR = '__build'
  env.GOPATH    = "${WORKSPACE}/${BUILD_DIR}"
  env.SRC_PATH  = "${env.GOPATH}/src/github.com/${REPO}"

  def govuk = load '/var/lib/jenkins/groovy_scripts/govuk_jenkinslib.groovy'

  try {
    stage("Checkout") {
       checkout scm
    }

    stage("Setup build environment") {
      // Clean GOPATH: Recursively delete everything in the current directory
      dir(env.GOPATH) {
        deleteDir()
      }

      // Create build path
      sh "mkdir -p ${env.SRC_PATH}"

      // Seed build path
      dir(env.WORKSPACE) {
        sh "/usr/bin/rsync -a ./ ${env.SRC_PATH} --exclude=$BUILD_DIR"
      }
    }

    stage("Build") {
      dir(env.SRC_PATH) {
        sh 'BINARY=$WORKSPACE/metadata-api make clean build'
      }
    }

    // Run tests
    wrap([$class: 'AnsiColorBuildWrapper']) {
      stage("Test") {
        dir(env.SRC_PATH) {
          sh 'BINARY=$WORKSPACE/metadata-api make test'
          // TODO: This appears to hang on CI.
          // sh '$WORKSPACE/metadata-api -version'
        }
      }
    }

    stage("Archive artefact") {
      archiveArtifacts 'metadata-api'
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
