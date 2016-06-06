#!/bin/bash
set -e

REPO=alphagov/metadata-api
export GOPATH=$PWD/gopath
GO_GITHUB_PATH=$GOPATH/src/github.com
BUILD_PATH=$GO_GITHUB_PATH/$REPO

rm -rf $GOPATH && mkdir -p $GOPATH/bin $BUILD_PATH

rsync -a ./ $BUILD_PATH --exclude=gopath

cd $BUILD_PATH && make
cp ./metadata-api $WORKSPACE/metadata-api
