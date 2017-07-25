#!/usr/bin/env bash

version=0.2.0

cd ~
nplh_base="$(pwd)/.nplh"

abspath() {
    if [ -d "$1" ]; then
        # dir
        (cd "$1"; pwd)
    elif [ -f "$1" ]; then
        # file
        if [[ $1 == */* ]]; then
            echo "$(cd "${1%/*}"; pwd)/${1##*/}"
        else
            echo "$(pwd)/$1"
        fi
    fi
}

try_curl() {
  command -v curl > /dev/null &&
  if [[ $1 =~ tgz$ ]]; then
    curl -fL $1 | tar -xzf -
  else
    local temp=${TMPDIR:-/tmp}/nplh.zip
    curl -fLo "$temp" $1 && unzip -o "$temp" && rm -f "$temp"
  fi
}

try_wget() {
  command -v wget > /dev/null &&
  if [[ $1 =~ tgz$ ]]; then
    wget -O - $1 | tar -xzf -
  else
    local temp=${TMPDIR:-/tmp}/nplh.zip
    wget -O "$temp" $1 && unzip -o "$temp" && rm -f "$temp"
  fi
}

download() {
  echo "Downloading nplh ..."
  mkdir -p "$nplh_base"/bin && cd "$nplh_base"/bin
  if [ $? -ne 0 ]; then
    binary_error="Failed to create bin directory"
    return
  fi

  url=https://github.com/nplh/nplh/releases/download/$version/$1

  echo downloading $url
  set -o pipefail
  if ! (try_curl $url || try_wget $url); then
    set +o pipefail
    binary_error="Failed to download with curl and wget"
    return
  fi
  set +o pipefail

  if [ ! -f nplh ]; then
    binary_error="Failed to download ${1}"
    return
  fi
  
  chmod +x nplh

  if [ "$(readlink /usr/bin/nplh)" != "$(abspath ~/.nplh/bin/nplh)" ]; then
    echo Linking binary
    sudo ln -s ~/.nplh/bin/nplh /usr/bin/nplh
  fi
}

archi=$(uname -sm)
binary_available=1
binary_error=""
case "$archi" in
  Darwin\ *64)   download nplh-$version-darwin-${binary_arch:-amd64}.tgz  ;;
  Darwin\ *86)   download nplh-$version-darwin-${binary_arch:-386}.tgz    ;;
  Linux\ *64)    download nplh-$version-linux-${binary_arch:-amd64}.tgz   ;;
  Linux\ *86)    download nplh-$version-linux-${binary_arch:-386}.tgz     ;;
  Linux\ armv5*) download nplh-$version-linux-${binary_arch:-arm5}.tgz    ;;
  Linux\ armv6*) download nplh-$version-linux-${binary_arch:-arm6}.tgz    ;;
  Linux\ armv7*) download nplh-$version-linux-${binary_arch:-arm7}.tgz    ;;
  Linux\ armv8*) download nplh-$version-linux-${binary_arch:-arm8}.tgz    ;;
  FreeBSD\ *64)  download nplh-$version-freebsd-${binary_arch:-amd64}.tgz ;;
  FreeBSD\ *86)  download nplh-$version-freebsd-${binary_arch:-386}.tgz   ;;
  OpenBSD\ *64)  download nplh-$version-openbsd-${binary_arch:-amd64}.tgz ;;
  OpenBSD\ *86)  download nplh-$version-openbsd-${binary_arch:-386}.tgz   ;;
  CYGWIN*\ *64)  download nplh-$version-windows-${binary_arch:-amd64}.zip ;;
  *)             binary_available=0 binary_error=1 ;;
esac

cd "$nplh_base"
if [ -n "$binary_error" ]; then
  if [ $binary_available -eq 0 ]; then
    echo "No prebuilt binary for $archi ..."
  else
    echo "  - $binary_error !!!"
  fi
  if command -v go > /dev/null; then
    echo -n "Building binary (go get -u github.com/nplh/nplh) ... "
    if [ -z "${GOPATH-}" ]; then
      export GOPATH="${TMPDIR:-/tmp}/nplh-gopath"
      mkdir -p "$GOPATH"
    fi
    if go get -u github.com/nplh/nplh; then
      echo "OK"
      cp "$GOPATH/bin/nplh" "$nplh_base/bin/"
    else
      echo "Failed to build binary. Installation failed."
      exit 1
    fi
  else
    echo "go executable not found. Installation failed."
    exit 1
  fi
fi
