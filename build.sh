#!/usr/bin/env bash

# clear build dir
rm -rf build
mkdir build

if git describe --exact-match HEAD > /dev/null 2>&1; then
  version=$(git describe --exact-match HEAD)
  build_release=true
else
  version=$(git describe HEAD)
  build_release=false
fi

build() {
  if [ "$1" == "arm" ]; then
    os=""
    arch="arm"
    arm="$2"
    desc="linux-arm$2"
  elif [ "$1" == "arm8" ]; then
    os=""
    arch="arm64"
    arm=""
    desc="linux-arm8"
  else
    os=$1
    arch=$2
    arm=""
    desc="$os-$arch"
  fi
  filename="build/nplh-$version-$desc"

  echo $filename

  echo "-  compiling"
  GOOS="$os" GOARCH="$arch" GOARM="$arm" CGO_ENABLED=0 go build

  if [ "$os" == "windows" ]; then
    echo "-  zipping"
    zip -q $filename.zip nplh
  else
    echo "-  tarring"
    tar -czf $filename.tgz nplh
  fi
  rm nplh
}

if [ "$build_release" == true ]; then
  build linux amd64
  build linux 386

  build windows amd64
  build windows 386

  build darwin amd64
  build darwin 386

  build arm 7
  build arm8
else
  build linux amd64
fi
