#!/bin/sh

set -e

function cleanRebuild {
    rm -rf bin/
    mkdir -p bin/
}

function buildServerRelease {
    echo "Building upspeak server release"
    CGO_ENABLED=0 GOOS=linux go build -o bin/upspeak -a -ldflags '-extldflags "-static"' .
}

function buildServerDev {
    echo "Building upspeak server dev"
    go build -o bin/upspeak .
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
