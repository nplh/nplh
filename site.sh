#!/usr/bin/env sh

# clear build dir
rm -rf build
mkdir build

version=1.0.0

redirects=""

redirect() {
  echo "$redirects$1  $2" >> build/_redirects
}

cp install.sh  build/index.html
cp install.sh  build/dl

# redirect / https://github.com/nplh/nplh
