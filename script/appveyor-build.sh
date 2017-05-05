#!/bin/sh
set -ex

export PATH=/c/msys64/mingw64/bin:/c/msys64/usr/bin:/c/Go/bin:/c/gopath/go/bin:$PATH
export GOROOT=/c/Go/
export GOPATH=/c/gopath

go get -d github.com/git-time-metric/git2go
cd /c/gopath/src/github.com/git-time-metric/git2go
git checkout v25
git submodule update --init

export BUILD="$PWD/vendor/libgit2/build"
export PCFILE="$BUILD/libgit2.pc"

FLAGS=$(pkg-config --static --libs $PCFILE) || exit 1
if [[ "$OSTYPE" == "msys" ]]; then 
  if [[ ! "lws2_32" == *"${FLAGS}"* ]]; then
    FLAGS="${FLAGS} -lws2_32"
  fi
fi
export CGO_LDFLAGS="$BUILD/libgit2.a -L$BUILD ${FLAGS}"
export CGO_CFLAGS="-I$PWD/vendor/libgit2/include"
go install ./...

cd /c/gopath/src/github.com/git-time-metric/gtm
go get -t -v ./...
go test --tags static -v ./...
if [[ "${APPVEYOR_REPO_TAG}" = true ]]; then
    go build --tags static -v -ldflags "-X main.Version=${APPVEYOR_REPO_TAG_NAME}"
    tar -zcf gtm.${APPVEYOR_REPO_TAG_NAME}.windows.tar.gz gtm.exe
else
    timestamp=$(date +%s)
    go build --tags static -v -ldflags "-X main.Version=developer-build-$timestamp"
    tar -zcf "gtm.developer-build-$timestamp.windows.tar.gz" gtm.exe
fi
