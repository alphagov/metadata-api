#!/bin/bash
set -x

VENV_PATH="${HOME}/venv/${JOB_NAME}"

[ -x ${VENV_PATH}/bin/pip ] || virtualenv ${VENV_PATH}
. ${VENV_PATH}/bin/activate

pip install -q ghtools

REPO=alphagov/metadata-api
export GOPATH=$PWD/gopath
GO_GITHUB_PATH=$GOPATH/src/github.com
BUILD_PATH=$GO_GITHUB_PATH/$REPO

rm -rf $GOPATH && mkdir -p $GOPATH/bin $BUILD_PATH

rsync -a ./ $BUILD_PATH --exclude=gopath

gh-status "$REPO" "$GIT_COMMIT" pending -d "\"Build #${BUILD_NUMBER} is running on Jenkins\"" -u "$BUILD_URL" >/dev/null

if cd $BUILD_PATH && make; then
  gh-status "$REPO" "$GIT_COMMIT" success -d "\"Build #${BUILD_NUMBER} succeeded on Jenkins\"" -u "$BUILD_URL" >/dev/null
  cp ./metadata-api $WORKSPACE/metadata-api
  exit 0
else
  gh-status "$REPO" "$GIT_COMMIT" failure -d "\"Build #${BUILD_NUMBER} failed on Jenkins\"" -u "$BUILD_URL" >/dev/null
  exit 1
fi
