#!/usr/bin/env sh

# clear build dir
rm -rf build
mkdir build

version=0.1.1

redirects=""

redirect() {
  echo "$redirects$1  $2" >> build/_redirects
}

cp install.sh  build/dl

redirect / https://github.com/nplh/nplh
