#!/bin/sh

set -e
set -x

# disable CGO for cross-compiling
export CGO_ENABLED=0

# compile for linux multi-arch
GOOS=linux GOARCH=amd64 go build -o release/linux/amd64/drone-gc   github.com/drone/drone-gc
GOOS=linux GOARCH=arm64 go build -o release/linux/arm64/drone-gc   github.com/drone/drone-gc
GOOS=linux GOARCH=arm   go build -o release/linux/arm/drone-gc     github.com/drone/drone-gc
