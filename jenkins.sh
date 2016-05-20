REPO=alphagov/metadata-api
export GOPATH=$PWD/gopath
GO_GITHUB_PATH=$GOPATH/src/github.com
BUILD_PATH=$GO_GITHUB_PATH/$REPO

rm -rf $GOPATH && mkdir -p $GOPATH/bin $BUILD_PATH

rsync -a ./ $BUILD_PATH --exclude=gopath

if cd $BUILD_PATH && make; then
  cp ./metadata-api $WORKSPACE/metadata-api
  exit 0
else
  exit 1
fi
