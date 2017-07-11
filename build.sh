# apk update
# apk upgrade
# apk add git curl
go get
go build ./nplh.go
mkdir build
mv nplh build
ls build
curl --request POST --header "PRIVATE-TOKEN: $APIKEY" --form "file=nplh" https://gitlab.example.com/api/v3/projects/nplh%2Fnplh/uploads
