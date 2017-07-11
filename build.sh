#!/usr/bin/env sh

apk update
apk upgrade
apk add git curl jq
go get
go build ./nplh.go
mkdir build
mv nplh build

if git describe --exact-match HEAD; then
  version=$(git describe --exact-match HEAD)
  echo "Uploading bin for $version"
  release_binary="$(curl \
    --request POST \
    --header "PRIVATE-TOKEN: $APIKEY" \
    --form "file=@build/nplh" \
    https://gitlab.com/api/v3/projects/nplh%2Fnplh/uploads | \
    jq -r '.markdown')"

  echo $release_binary

  curl \
    --request POST \
    --header "PRIVATE-TOKEN: $APIKEY" \
    --data $release_binary \
    https://gitlab.com/api/v3/projects/nplh%2Fnplh/repository/tags/$version/release | \
    jq
else
  echo "Not a tag; not uploading"
fi
