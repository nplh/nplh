#!/usr/bin/env sh

apk update
apk upgrade
apk add git curl jq gcc
go get
GOOS=linux GOARCH=amd64 go build -tags netgo -installsuffix netgo ./nplh.go
file nplh
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
    --header "Content-Type: application/json" \
    --data "{'description': $release_binary}" \
    https://gitlab.com/api/v3/projects/nplh%2Fnplh/repository/tags/$version/release

else
  echo "Not a tag; not uploading"
fi
