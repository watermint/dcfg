#!/usr/bin/env bash

BUILD=`pwd`/out
DIST=/dist
if [ "$BUILD_ID"x = ""x ]; then
  BUILD_ID=0
fi
APP_VERSION=`cat $PROJECT_ROOT/version`.$BUILD_ID

echo building: $APP_VERSION
echo UID: `id`

cd $PROJECT_ROOT
glide install

echo Testing...
go test  $(glide novendor)
if [ x"$?" != x"0" ]; then
  echo Test failed: $?
  exit 1
fi

LD_FLAGS="-X main.AppVersion=$APP_VERSION"

GOOS=windows GOARCH=386   go build --ldflags "$LD_FLAGS" -o $BUILD/dcfg-$APP_VERSION-win.exe github.com/watermint/dcfg
GOOS=linux   GOARCH=386   go build --ldflags "$LD_FLAGS" -o $BUILD/dcfg-$APP_VERSION-linux   github.com/watermint/dcfg
GOOS=darwin  GOARCH=amd64 go build --ldflags "$LD_FLAGS" -o $BUILD/dcfg-$APP_VERSION-darwin  github.com/watermint/dcfg

cd $BUILD
zip -9 -r $DIST/dcfg-$APP_VERSION.zip .
