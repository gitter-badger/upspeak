#!/bin/bash

set -e

function cleanRebuild {
    rm -rf bin/
    mkdir -p bin/
}

function buildServerRelease {
    echo "Building upspeak server release"
    export CGO_ENABLED=0
    export GOOS=linux
    go build -o bin/upspeak-rig -a -ldflags '-extldflags "-static"' .
    go build -o bin/upspeak-migrate -a -ldflags '-extldflags "-static"' ./tools/migrate.go
}

function buildServerDev {
    echo "Building upspeak server dev"
    go build -o bin/upspeak-rig .
    go build -o bin/upspeak-migrate ./tools/migrate.go
}

function buildDev {
    cleanRebuild
    buildServerDev
}

function buildRelease {
    cleanRebuild
    buildServerRelease
}

case "$1" in
    "release")
    buildRelease
    ;;
    "dev")
    buildDev
    ;;
    *)
    echo "Usage: ./build.sh [dev|release|help]"
    echo "Binaries are put in './bin/'"
    ;;
esac
