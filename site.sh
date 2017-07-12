#!/usr/bin/env sh

rm -rf build
mkdir build

redirects=""

redirect() {
  redirects=$redirects$1"   "$2$'\n'
}

latest=$(git describe)

redirect /dl $latest

echo "$redirects"
echo "$redirects" > build/_redirects
