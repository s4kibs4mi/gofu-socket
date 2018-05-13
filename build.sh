#! /bin/sh
set -e

export GOPATH=$HOME/go/

if ! [ -x "$(command -v go)" ]; then
    echo "go is not installed"
    exit
fi
if ! [ -x "$(command -v git)" ]; then
    echo "git is not installed"
    exit
fi
if [ -z "$GOPATH" ]; then
    echo "set GOPATH"
    exit
fi

export GOARCH="amd64"
export GOOS="linux"
export CGO_ENABLED=0

go get .
go build -o gofu-socket -v .

docker build . -t gofu-socket:latest
docker push gofu-socket:latest

rm gofu-socket
